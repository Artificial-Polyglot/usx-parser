package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Artificial-Polyglot/io-lib/dbio"
	"github.com/Artificial-Polyglot/io-lib/safe"
)

type stack []string

type script struct {
	BookId      string
	ChapterNum  int
	VerseStr    string
	VerseEnd    string
	VerseNum    int
	ScriptNum   string
	UsfmStyle   string
	ScriptTexts []string
}

type titleDesc struct {
	heading string
	title   []string
}

type USXParser struct {
	db  *sql.DB
	log *dbio.Logger
}

var hasStyle = map[string]bool{
	`book`: true, `para`: true, `char`: true, `cell`: true,
	`ms`: true, `note`: true, `sidebar`: true, `figure`: true,
}

func (p *USXParser) parseUSXFiles(files []string) error {
	var allRecords []script

	for _, filename := range files {
		records, t, err := p.decodeUSX(filename)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", filepath.Base(filename), err)
		}
		records = p.addChapterHeading(records, t)
		records = p.correctScriptNum(records)
		allRecords = append(allRecords, records...)
	}

	if len(allRecords) == 0 {
		return fmt.Errorf("no records parsed from %d files", len(files))
	}
	return p.insertScripts(allRecords)
}

func (p *USXParser) decodeUSX(filename string) ([]script, titleDesc, error) {
	var records []script
	var titles titleDesc

	xmlFile, err := os.Open(filename)
	if err != nil {
		return records, titles, fmt.Errorf("opening %s: %w", filename, err)
	}
	defer xmlFile.Close()

	var stk stack
	var rec script
	var tagName string
	var chapterNum = 1
	var scriptNum = 0
	var verseNum int
	var verseStr = `0`
	var usfmStyle string

	decoder := xml.NewDecoder(xmlFile)
	for {
		var token xml.Token
		token, err = decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return records, titles, fmt.Errorf("parsing XML: %w", err)
		}
		switch se := token.(type) {
		case xml.StartElement:
			tagName = se.Name.Local
			if tagName == `book` {
				rec.BookId = p.findAttr(se, `code`)
			} else if tagName == `chapter` {
				chapterNum = p.findIntAttr(se, `number`)
			} else if tagName == `verse` {
				verseStr = p.findAttr(se, `number`)
				verseNum = p.findIntAttr(se, `number`)
			}
			if hasStyle[tagName] {
				usfmStyle = tagName + `.` + p.findAttr(se, `style`)
				if p.include(usfmStyle) {
					stk = stk.push(usfmStyle)
				} else {
					err = decoder.Skip()
					if err != nil {
						return records, titles, fmt.Errorf("skipping element: %w", err)
					}
				}
			}
		case xml.CharData:
			text := string(se)
			if len(strings.TrimSpace(text)) > 0 {
				if strings.Contains(text, "{") || strings.Contains(text, "}") {
					text = strings.Replace(text, `{`, `(`, -1)
					text = strings.Replace(text, `}`, `)`, -1)
				}
				if usfmStyle == `para.h` {
					titles.heading = text
				} else if usfmStyle == `para.mt` || usfmStyle == `para.mt1` ||
					usfmStyle == `para.mt2` || usfmStyle == `para.mt3` {
					titles.title = append(titles.title, text)
				} else {
					rec.ScriptTexts = append(rec.ScriptTexts, text)
				}
			}
		case xml.EndElement:
			if hasStyle[se.Name.Local] {
				stk, usfmStyle = stk.pop()
			}
		}
		if chapterNum != rec.ChapterNum || verseNum != rec.VerseNum {
			if rec.BookId != `` && len(rec.ScriptTexts) > 0 {
				records = append(records, rec)
			}
			scriptNum++
			if chapterNum != rec.ChapterNum {
				scriptNum = 1
			}
			rec = script{
				BookId:     rec.BookId,
				ChapterNum: chapterNum,
				ScriptNum:  strconv.Itoa(scriptNum),
				VerseNum:   verseNum,
				VerseStr:   verseStr,
				UsfmStyle:  usfmStyle,
			}
		}
	}
	if rec.BookId != `` && len(rec.ScriptTexts) > 0 {
		records = append(records, rec)
	}
	return records, titles, nil
}

func (p *USXParser) addChapterHeading(records []script, titles titleDesc) []script {
	results := make([]script, 0, len(records)+300)
	if len(records) == 0 {
		return results
	}
	rec := records[0]
	rec.VerseStr = `0`
	rec.VerseNum = 0
	rec.UsfmStyle = `para.mt`
	rec.ScriptTexts = []string{strings.Join(titles.title, " ")}
	results = append(results, rec)
	lastChapter := 1
	for _, rec = range records {
		if lastChapter != rec.ChapterNum {
			lastChapter = rec.ChapterNum
			rec2 := rec
			rec2.VerseStr = `0`
			rec2.VerseNum = 0
			rec2.UsfmStyle = `para.h`
			rec2.ScriptTexts = []string{titles.heading + " " + strconv.Itoa(rec.ChapterNum)}
			results = append(results, rec2)
		}
		results = append(results, rec)
	}
	return results
}

func (p *USXParser) correctScriptNum(records []script) []script {
	results := make([]script, 0, len(records))
	scriptNum := 0
	lastChapter := 0
	for _, rec := range records {
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			scriptNum = 0
		}
		scriptNum++
		rec.ScriptNum = strconv.Itoa(scriptNum)
		results = append(results, rec)
	}
	return results
}

func (p *USXParser) findAttr(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ``
}

func (p *USXParser) findIntAttr(se xml.StartElement, name string) int {
	val := p.findAttr(se, name)
	if val == `` {
		return 0
	}
	return safe.SafeVerseNum(val)
}

func (p *USXParser) insertScripts(records []script) error {
	records = parseVerseStr(records)

	query := `INSERT INTO scripts(dataset_id, book_id, chapter_num, chapter_end, audio_file, script_num, usfm_style,
		person, actor, verse_num, verse_str, verse_end, script_text)
		VALUES (1,?,?,0,'',?,?,'','',?,?,?,?)`

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, rec := range records {
		scriptNum := dbio.ZeroFill(rec.ScriptNum, 5)
		text := safe.SafeStringJoin(rec.ScriptTexts)
		_, err = stmt.Exec(rec.BookId, rec.ChapterNum, scriptNum, rec.UsfmStyle,
			rec.VerseNum, rec.VerseStr, rec.VerseEnd, text)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("inserting script: %w", err)
		}
	}
	return tx.Commit()
}

func parseVerseStr(records []script) []script {
	for i := range records {
		parts := strings.Split(records[i].VerseStr, "-")
		if len(parts) > 1 {
			records[i].VerseStr = parts[0]
			records[i].VerseEnd = parts[len(parts)-1]
		}
	}
	return records
}

func (s stack) push(v string) stack {
	return append(s, v)
}

func (s stack) pop() (stack, string) {
	l := len(s)
	if l < 1 {
		return s, ""
	}
	return s[:l-1], s[l-1]
}
