package elf

/* Executable and Linking Format */

/*
	1. A relocatable file holds code and data suitable for linking with other object files to create an executable or a shared object file.
	2. An executable file holds a program suitable for execution; the file specifies how exec(BA_OS) creates
	a program’s process image.
	3. A shared object file holds code and data suitable for linking in two contexts. First, the link editor [see ld(SD_CMD)] may process it with other relocatable and shared object files to create another object file. Second, the dynamic linker combines it with an executable file and other shared objects to create a process image.

	An ELF header resides at the beginning and holds a ‘‘road map’’ describing the file’s organization. Sections hold the bulk of object file information for the linking view: instructions, data, symbol table, relocation information, and so on.

	A program header table, if present, tells the system how to create a process image. Files used to build a process image (execute a program) must have a program header table; relocatable files do not need one. A section header table contains information describing the file’s sections. Every section has an entry in the table; each entry gives information such as the section name, the section size, etc. Files used during linking must have a section header table; other object files may or may not have one.

	 The object file format supports various processors with 8-bit bytes and 32-bit architectures. Nevertheless, it is intended to be extensible to larger (or smaller) architectures. Object files therefore represent some control data with a machine-independent format, making it possible to identify object files and interpret their contents in a common way. Remaining data in an object file use the encoding of the target processor, regardless of the machine on which the file was created.

ELF Header.
	 #define EI_NIDENT 16
	 typedef struct {
		 unsigned char e_ident[EI_NIDENT];
		 Elf32_Half e_type;             // ET_NONE/ET_REL/ET_EXEC/ET_DYN/ET_CORE/ET_LOPROC/ET_HIPROC
		 Elf32_Half e_machine;          // EM_NONE/EM_M32/EM_SPARC/EM_386/EM_68K/EM_88K/EM_860/EM_MIPS
		 Elf32_Word e_version;          // EV_NONE/EV_CURRENT
		 Elf32_Addr e_entry;            // This member gives the virtual address to which the system first transfers control, thus starting the process, If the file has no associated entry point, this member holds zero.
		 Elf32_Off e_phoff;             // This member Holds the program header table's file offset in bytes. If the file has no program header table, this member holds zero.
		 Elf32_Off e_shoff;             // This member holds the section header table's file offset in bytes. If the file has no section header table, this member holds zero.
		 Elf32_Word e_flags;
		 Elf32_Half e_ehsize;           // This member holds the ELF header's size in bytes.
		 Elf32_Half e_phentsize;        // This member holds the size in bytes of one entry in the file's program header table; all entries are the same size.
		 Elf32_Half e_phnum;            // This member holds the number of entries in the program header table. Thus the product of e_phentsize and e_phnum gives the table's size in bytes. If a file has no program header table, e_phnum table holds the value zero.
		 Elf32_Half e_shentsize;        // This member holds a section header's size in bytes. A section header is one entry in the section header table; all entries are the same size.
		 Elf32_Half e_shnum;            // This member hold the number of entries in the section header table. Thus the product of e_shentsize and e_shnum gives the section header table's size in bytes. If a file has no section header table, e_shnum holds the value zero.
		 Elf32_Half e_shstrndx;         // This member holds the section header table index of the entry associated with the section name string table. If the file has no section nmae string table, this member holds the value SHN_UNDEF.
	 } Elf32_Ehdr;

ELF Identification
	 As mentioned above, ELF provides an object file framework to support multiple processors, multiple data
	 encodings, and multiple classes of machines. To support this object file family, the initial bytes of the file
	 specify how to interpret the file, independent of the processor on which the inquiry is made and independent of the file’s remaining contents.
	 e_ident[ ] Identification Indexes
	 name            value           purpose
	 ___________________________________________
	 EI_MAG0          0              File identification
	 EI_MAG1          1              File identification
	 EI_MAG2          2              File identification
	 EI_MAG3          3              File identification
	 EI_CLASS         4              File class
	 EI_DATA          5              Data encoding
	 EI_VERSION       6              File version
	 EI_PAD           7              Start of padding bytes
	 EI_NIDENT        16             Size of e_ident[]
	 ___________________________________________

Sections
	An object file’s section header table lets one locate all the file’s sections.
	The section header table is an array of Elf32_Shdr structures as described below. A section header table index is a subscript into this array. The ELF header’s e_shoff member gives the byte offset from the beginning of the file to the section header table; e_shnum tells how many entries the section header table contains; e_shentsize gives the size in bytes of each entry.
	Some section header table indexes are reserved; an object file will not have sections for these special indexes.
	Sections contain all information in an object file, except the ELF header, the program header table, and the section header table. Moreover, object files’ sections satisfy several condition:
		1. Every section in an object file has exactly one section header describing it. Section headers may exist that do not have a section.
		2. Each section occupies one contiguous (possibly empty) sequence of bytes within a file.
		3. Sections in a file may not overlap. No byte in a file resides in more than one section.
		4. An object file may have inactive space. The various headers and the sections might not ‘‘cover’’ every byte in an object file. The contents of the inactive data are unspecified.

	Section Header
	typedef struct {
		Elf32_Word sh_name;
		Elf32_Word sh_type;
		Elf32_Word sh_flags;
		Elf32_Addr sh_addr;
		Elf32_Off sh_offset;
		Elf32_Word sh_size;
		Elf32_Word sh_link;
		Elf32_Word sh_info;
		Elf32_Word sh_addralign;
		Elf32_Word sh_entsize;
	} Elf32_Shdr;

	Section Types, sh_type
	Name Value
	_ _____________________________
	SHT_NULL 			0
	SHT_PROGBITS 		1
	SHT_SYMTAB 			2
	SHT_STRTAB 			3
	SHT_RELA 			4
	SHT_HASH 			5
	SHT_DYNAMIC 		6
	SHT_NOTE 			7
	SHT_NOBITS 			8
	SHT_REL 			9
	SHT_SHLIB 			10
	SHT_DYNSYM 			11
	SHT_LOPROC 			0x70000000
	SHT_HIPROC 			0x7fffffff
	SHT_LOUSER 			0x80000000
	SHT_HIUSER 			0xffffffff
	_ _____________________________ 
*/
type ELF struct {
}
