package memory

// Header table
const (
	HVersion                       Address = 0x00
	HFlags                         Address = 0x01
	HHighMemoryBegin               Address = 0x04
	HMainRoutine                   Address = 0x06
	HDictionary                    Address = 0x08
	HObjectTable                   Address = 0x0A
	HGlobalVariables               Address = 0x0C
	HStaticMemoryBegin             Address = 0x0E
	HFlags2                        Address = 0x10
	HAbbreviationsTable            Address = 0x18 // v2+
	HFileLength                    Address = 0x1A // v3+ (not all v3)
	HFileChecksum                  Address = 0x1C // v3+ (not all v3)
	HInterpreterNumber             Address = 0x1E // v4+
	HInterpreterVersion            Address = 0x1F // v4+
	HScreenHeightLines             Address = 0x20 // v4+
	HScreenWidthWidthChars         Address = 0x21 // v4+
	HScreenWidthUnits              Address = 0x22 // v5+
	HScreenHeightUnits             Address = 0x24 // v5+
	HFontWidthUnitsV5              Address = 0x26 // v5 ONLY
	HFontHeightUnits               Address = 0x26 // v6+
	HFontHeightUnitsV5             Address = 0x27 // v5 ONLY
	HFontWidthUnits                Address = 0x27 // v6+
	HRoutinesOffset                Address = 0x28 // v6+
	HStaticStringsOffset           Address = 0x2A // v6+
	HDefaultBackgroundColour       Address = 0x2C // v5+
	HDefaultForegroundColour       Address = 0x2D // v5+
	HTerminatingCharactersTable    Address = 0x2E // v5+
	HTotalWidthSentToOutputStream3 Address = 0x30 // v6+
	HStandardRevisionNumber        Address = 0x32
	HAlphabetTable                 Address = 0x34 // v5+
	HHeaderExtensionTable          Address = 0x36 // v5+
)
