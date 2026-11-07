package src

type Language string

const (
	LangEN Language = "en"
	LangES Language = "es"
)

func (l Language) Next() Language {
	if l == LangEN {
		return LangES
	}
	return LangEN
}

type Dict struct {
	AppTitle        string
	AppSubtitle     string
	Scanning        string
	ScanDone        string
	Calculating     string
	ParentDir       string
	CurrentDir      string
	EmptyDir        string
	HeaderRank      string
	HeaderSize      string
	HeaderUsage     string
	HeaderName      string
	SortLabel       string
	SortSizeDesc    string
	SortSizeAsc     string
	SortNameAsc     string
	SortNameDesc    string
	SortCountDesc   string
	HiddenShown     string
	HiddenHidden    string
	PreviewTitle    string
	TopSubdirs      string
	TopFiles        string
	TotalContained  string
	TotalFiles      string
	TotalDirs       string
	Modified        string
	Permissions     string
	ShareOfParent   string
	NoFiles         string
	NoSubdirs       string
	BinaryFile      string
	PreviewLines    string
	RemainingLines  string
	HelpFooter      string
	AboutTitle      string
	AboutSubtitle   string
	AboutVersion    string
	AboutCreator    string
	AboutRepo       string
	AboutLicense    string
	AboutLicenseVal string
	AboutYear       string
	AboutBuiltWith  string
	AboutReturnHint string
	ErrAccess       string
	ErrReadDir      string
}

var dictionaries = map[Language]Dict{
	LangEN: {
		AppTitle:        "bigcrock",
		AppSubtitle:     "Disk usage tree analyzer and directory ranker",
		Scanning:        "Scanning... (%d/%d done)",
		ScanDone:        "Ready (%d items, %s total)",
		Calculating:     "Calculating...",
		ParentDir:       "Parent",
		CurrentDir:      "Current Directory",
		EmptyDir:        "(Empty directory)",
		HeaderRank:      "RANK",
		HeaderSize:      "SIZE",
		HeaderUsage:     "USAGE",
		HeaderName:      "DIRECTORY / FILE",
		SortLabel:       "Sort:",
		SortSizeDesc:    "Size (Largest First)",
		SortSizeAsc:     "Size (Smallest First)",
		SortNameAsc:     "Name (A-Z)",
		SortNameDesc:    "Name (Z-A)",
		SortCountDesc:   "Items Count",
		HiddenShown:     "Hidden: ON",
		HiddenHidden:    "Hidden: OFF",
		PreviewTitle:    "DETAILS & BREAKDOWN",
		TopSubdirs:      "HEAVIEST SUBDIRECTORIES",
		TopFiles:        "LARGEST FILES",
		TotalContained:  "Total Size:",
		TotalFiles:      "Files:",
		TotalDirs:       "Subdirectories:",
		Modified:        "Modified:",
		Permissions:     "Permissions:",
		ShareOfParent:   "Share of directory:",
		NoFiles:         "(No regular files)",
		NoSubdirs:       "(No subdirectories)",
		BinaryFile:      "Binary or Media File (%s)",
		PreviewLines:    "Content Preview (first %d lines):",
		RemainingLines:  "... [%d remaining lines]",
		HelpFooter:      "[k/j] Move  [l/Enter] Open  [h] Up  [s] Sort  [.] Hidden  [r] Rescan  [Tab/t] Lang: EN  [?] About  [q] Quit",
		AboutTitle:      "bigcrock",
		AboutSubtitle:   "Disk usage tree analyzer and directory ranker, written in Go",
		AboutVersion:    "Version:",
		AboutCreator:    "Creator:",
		AboutRepo:       "Repository:",
		AboutLicense:    "License:",
		AboutLicenseVal: "MIT License",
		AboutYear:       "Year:",
		AboutBuiltWith:  "Made with:",
		AboutReturnHint: "Press any key to return",
		ErrAccess:       "Error: cannot access %q: %v",
		ErrReadDir:      "Error reading directory: %v",
	},
	LangES: {
		AppTitle:        "bigcrock",
		AppSubtitle:     "Analizador de espacio en disco y ranking de carpetas",
		Scanning:        "Analizando... (%d/%d listos)",
		ScanDone:        "Listo (%d elementos, %s total)",
		Calculating:     "Calculando...",
		ParentDir:       "Padre",
		CurrentDir:      "Directorio Actual",
		EmptyDir:        "(Directorio vacío)",
		HeaderRank:      "RANGO",
		HeaderSize:      "PESO",
		HeaderUsage:     "USO",
		HeaderName:      "DIRECTORIO / ARCHIVO",
		SortLabel:       "Orden:",
		SortSizeDesc:    "Peso (Mayor a Menor)",
		SortSizeAsc:     "Peso (Menor a Mayor)",
		SortNameAsc:     "Nombre (A-Z)",
		SortNameDesc:    "Nombre (Z-A)",
		SortCountDesc:   "Cantidad de Elementos",
		HiddenShown:     "Ocultos: SI",
		HiddenHidden:    "Ocultos: NO",
		PreviewTitle:    "DETALLES Y DESGLOSE",
		TopSubdirs:      "SUBDIRECTORIOS MAS PESADOS",
		TopFiles:        "ARCHIVOS MAS PESADOS",
		TotalContained:  "Tamaño Total:",
		TotalFiles:      "Archivos:",
		TotalDirs:       "Carpetas:",
		Modified:        "Modificado:",
		Permissions:     "Permisos:",
		ShareOfParent:   "Proporcion del directorio:",
		NoFiles:         "(Sin archivos regulares)",
		NoSubdirs:       "(Sin subdirectorios)",
		BinaryFile:      "Archivo Binario o Multimedia (%s)",
		PreviewLines:    "Vista Previa (primeras %d lineas):",
		RemainingLines:  "... [%d lineas restantes]",
		HelpFooter:      "[k/j] Mover  [l/Enter] Abrir  [h] Subir  [s] Orden  [.] Ocultos  [r] Re-escanear  [Tab/t] Idioma: ES  [?] Acerca de  [q] Salir",
		AboutTitle:      "bigcrock",
		AboutSubtitle:   "Analizador de espacio en disco y ranking de carpetas, escrito en Go",
		AboutVersion:    "Versión:",
		AboutCreator:    "Creador:",
		AboutRepo:       "Repositorio:",
		AboutLicense:    "Licencia:",
		AboutLicenseVal: "Licencia MIT",
		AboutYear:       "Año:",
		AboutBuiltWith:  "Hecho con:",
		AboutReturnHint: "Pulsa cualquier tecla para volver",
		ErrAccess:       "Error: no se pudo acceder a %q: %v",
		ErrReadDir:      "Error al leer directorio: %v",
	},
}

func T(lang Language) Dict {
	if d, ok := dictionaries[lang]; ok {
		return d
	}
	return dictionaries[LangEN]
}
