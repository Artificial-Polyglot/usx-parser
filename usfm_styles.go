package main

func (p *USXParser) include(style string) bool {
	last := style[len(style)-1]
	if last >= '0' && last <= '9' {
		style = style[:len(style)-1]
	}
	answer, ok := usfm[style]
	if !ok {
		p.log.Warn("USFM map does not have entry:", style)
	}
	return answer
}

var usfm = map[string]bool{
	`book.id`:   false,
	`chapter.c`: false,
	`verse.v`:   false,
	// Identification
	`para.ide`:  false,
	`para.h`:    true,
	`para.toc`:  false,
	`para.toca`: false,
	`para.rem`:  false,
	`para.usfm`: false,
	// Introductions
	`para.imt`:  false,
	`para.is`:   false,
	`para.ip`:   false,
	`para.ipi`:  false,
	`para.im`:   false,
	`para.imi`:  false,
	`para.ipq`:  false,
	`para.imq`:  false,
	`para.ipr`:  false,
	`para.iq`:   false,
	`para.ib`:   false,
	`para.ili`:  false,
	`para.iot`:  false,
	`para.io`:   false,
	`para.iex`:  false,
	`para.imte`: false,
	`para.ie`:   false,
	// Titles & Headings
	`para.mt`:  true,
	`para.mte`: false,
	`para.cl`:  false,
	`para.cd`:  false,
	`para.ms`:  false,
	`para.mr`:  false,
	`para.s`:   false,
	`para.sr`:  false,
	`para.r`:   false,
	`para.d`:   false,
	`para.sp`:  false,
	`para.sd`:  false,
	// Paragraphs
	`para.p`:       true,
	`para.m`:       true,
	`para.po`:      true,
	`para.pr`:      true,
	`para.cls`:     true,
	`para.pmo`:     true,
	`para.pm`:      true,
	`para.pmc`:     true,
	`para.pmr`:     true,
	`para.pi`:      true,
	`para.mi`:      true,
	`para.pc`:      true,
	`para.ph`:      true,
	`para.lit`:     true,
	`para.nb`:      true,
	`para.pb`:      true,
	`para.cp`:      false,
	`para.restore`: false,
	// Poetry
	`para.q`:  true,
	`para.qr`: true,
	`para.qc`: true,
	`para.qa`: false,
	`para.qm`: true,
	`para.qd`: false,
	`para.b`:  false,
	// Lists
	`para.lh`:   true,
	`para.li`:   true,
	`para.lf`:   true,
	`para.lim`:  true,
	`para.litl`: true,
	// Table
	`row.tr`:   true,
	`cell.th`:  true,
	`cell.thr`: true,
	`cell.tc`:  true,
	`cell.tcr`: true,
	// Identification chars
	`char.va`: false,
	`char.vp`: false,
	`char.ca`: false,
	// Special Text
	`char.add`:   false,
	`char.addnp`: false,
	`char.bk`:    true,
	`char.dc`:    false,
	`char.ior`:   false,
	`char.iqt`:   false,
	`char.k`:     true,
	`char.litl`:  true,
	`char.nd`:    true,
	`char.ord`:   true,
	`char.pn`:    true,
	`char.png`:   true,
	`char.qac`:   true,
	`char.qs`:    true,
	`char.qt`:    true,
	`char.rq`:    false,
	`char.sig`:   true,
	`char.sls`:   true,
	`char.tl`:    true,
	`char.wj`:    true,
	// Character Styling
	`char.em`:   true,
	`char.bd`:   true,
	`char.bdit`: true,
	`char.it`:   true,
	`char.no`:   true,
	`char.sc`:   true,
	`char.sup`:  true,
	// Special Features
	`char.rb`:  true,
	`char.pro`: true,
	`char.w`:   true,
	`char.wg`:  true,
	`char.wh`:  true,
	`char.wa`:  true,
	`char.fig`: false,
	// Structured List Entries
	`char.lik`: true,
	`char.liv`: true,
	// Linking
	`char.jmp`: true,
	// Note
	`note.f`:      false,
	`note.fe`:     false,
	`note.ef`:     false,
	`note.x`:      false,
	`note.ex`:     false,
	`sidebar.esb`: false,
	`figure.fig`:  false,
	// MS Section
	`ms.qt`: false,
	`ms.ts`: false,
}
