package file

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetFilePath(ctx context.Context) string {
	// type OpenDialogOptions struct {
	// 	DefaultDirectory           string
	// 	DefaultFilename            string
	// 	Title                      string
	// 	Filters                    []FileFilter
	// 	ShowHiddenFiles            bool
	// 	CanCreateDirectories       bool
	// 	ResolvesAliases            bool
	// 	TreatPackagesAsDirectories bool
	// }
	options := runtime.OpenDialogOptions{
		DefaultDirectory:           "",
		DefaultFilename:            "",
		Title:                      "选择文件",
		Filters:                    []runtime.FileFilter{{DisplayName: "All Files (*.*)", Pattern: "*.*"}},
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	}
	filePath, err := runtime.OpenFileDialog(ctx, options)
	if err != nil {
		return ""
	}
	return filePath

}
func GetDirPath(ctx context.Context) string {
	options := runtime.OpenDialogOptions{
		// DefaultDirectory:           "",
		// DefaultFilename:            "",
		// Title:                      "要下载到哪个文件夹？",
		// Filters:                    []runtime.FileFilter{{DisplayName: "All Files (*.*)", Pattern: "*.*"}},
		// ShowHiddenFiles:            false,
		// CanCreateDirectories:       false,
		// ResolvesAliases:            false,
		// TreatPackagesAsDirectories: false,
		Title: "要下载到哪个文件夹？",
		Filters: []runtime.FileFilter{
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
		CanCreateDirectories: true,
	}
	filePath, err := runtime.OpenDirectoryDialog(ctx, options)
	if err != nil {
		return ""
	}
	return filePath
}
