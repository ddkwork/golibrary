package safeType

import (
	"testing"
)

func TestNewHexDump(t *testing.T) {
	NewHexDump(HexDumpString(dump))
	NewHexDump(HexDumpString(bugBuf))
}

func TestHexDumpToGoBytes(t *testing.T) {
	ss := `00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00 8E 6A 64 01 15 4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47 00`
	New(HexDumpString(ss))
	NewHexDump(HexDumpString(dump))
	NewHexDump(`00 00 00 1A 00 00 00 09 00 01 00 00 0B 00 00 00`)             // 16 byte header
	NewHexDump(`8E 6A 64 01`)                                                 // tag
	NewHexDump(`01`)                                                          //
	NewHexDump(`15`)                                                          // strinf type id
	NewHexDump(`4F 53 44 4B 5F 41 42 55 53 45 5F 52 45 50 4F 52 54 49 4E 47`) // OSDK_ABUSE_REPORTING
	NewHexDump(`00`)                                                          // string end
}

var bugBuf = `
08A73200 57 61 72 68 61 6D 6D 65 72 20 34 30 2C 30 30 30  Warhammer 40,000  
08A73210 20 44 61 77 6E 20 6F 66 20 57 61 72 20 49 49 20   Dawn of War II   
08A73220 52 65 74 72 69 62 75 74 69 6F 6E 20 2D 20 49 6D  Retribution - Im  
08A73230 70 65 72 69 61 6C 20 46 69 73 74 73 20 43 68 61  perial Fists Cha  
08A73240 70 74 65 72 20 50 61 63 6B 00 00 00 00 00 00 00  pter Pack.......  
08A73250 57 61 72 68 61 6D 6D 65 72 20 34 30 2C 30 30 30  Warhammer 40,000  
08A73260 3A 20 44 61 77 6E 20 6F 66 20 57 61 72 20 49 49  : Dawn of War II  
08A73270 20 2D 20 52 65 74 72 69 62 75 74 69 6F 6E 20 2D   - Retribution -  
08A73280 20 49 6D 70 65 72 69 61 6C 20 46 69 73 74 73 20   Imperial Fists   
08A73290 43 68 61 70 74 65 72 20 50 61 63 6B 00 00 00 00  Chapter Pack....  
08A732A0 50 65 6E 6E 79 20 41 72 63 61 64 65 20 41 64 76  Penny Arcade Adv  
08A732B0 65 6E 74 75 72 65 73 20 4F 6E 20 74 68 65 20 52  entures On the R  
08A732C0 61 69 6E 2D 53 6C 69 63 6B 20 50 72 65 63 69 70  ain-Slick Precip  

`

var dump = `
00000000  7e 15 00 80 0b 00 00 00  09 25 ce f7 3d 01 00 10  |~........%..=...|
00000010  01 10 00 08 ac 80 04 1a  b2 01 32 00 00 00 04 00  |..........2.....|
00000020  00 00 25 ce f7 3d 01 00  10 01 07 00 00 00 83 7c  |..%..=.........||
00000030  39 6a 97 2b a8 c0 00 00  00 00 00 a9 25 63 80 58  |9j.+........%c.X|
00000040  41 63 01 00 00 00 00 00  00 00 00 00 ad 93 9b 27  |Ac.............'|
00000050  eb b6 3d dc 1d 57 fe b2  d1 86 79 de a1 41 61 eb  |..=..W....y..Aa.|
00000060  04 70 81 ce 35 f5 28 6a  05 52 d9 7b 7d 6c f9 2e  |.p..5.(j.R.{}l..|
00000070  5c b9 5e 8a b6 a5 87 dc  da 25 03 0b 00 48 76 7b  |\.^......%...Hv{|
00000080  66 ba f9 0b 48 78 62 09  bf 88 be 49 de 09 36 52  |f...Hxb....I..6R|
00000090  57 42 8d 69 34 8b 80 ac  e9 0b 8f ef e1 dd a2 0b  |WB.i4...........|
000000a0  25 0c cf 26 f9 0f dc 30  df 21 46 8f b6 8d c2 56  |%..&...0.!F....V|
000000b0  78 88 ef 2a 97 8c 50 c7  e2 9b 42 6f 53 09 82 42  |x..*..P...BoS..B|
000000c0  cc d4 3e 57 b5 ef b4 23  2c 54 13 97 20 d1 cf f0  |..>W...#,T.. ...|
000000d0  a7 b2 98 85 d3 54                                 |.....T|
`
