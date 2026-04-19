package manifest

import (
	"fmt"
	"sort"

	libfossil "github.com/danmestas/libfossil/internal/fsltype"
	"github.com/danmestas/libfossil/internal/blob"
	"github.com/danmestas/libfossil/internal/content"
	"github.com/danmestas/libfossil/internal/deck"
	"github.com/danmestas/libfossil/internal/repo"
)

type FileEntry struct {
	Name string
	UUID string
	Perm string
}

func ListFiles(r *repo.Repo, rid libfossil.FslID) ([]FileEntry, error) {
	d, err := GetManifest(r, rid)
	if err != nil {
		return nil, fmt.Errorf("manifest.ListFiles: %w", err)
	}
	if d.B == "" {
		return fCardsToEntries(d.F), nil
	}
	baseRid, ok := blob.Exists(r.DB(), d.B)
	if !ok {
		return nil, fmt.Errorf("manifest.ListFiles: baseline %s not found", d.B)
	}
	baseData, err := content.Expand(r.DB(), baseRid)
	if err != nil {
		return nil, fmt.Errorf("manifest.ListFiles: expand baseline: %w", err)
	}
	baseDeck, err := deck.Parse(baseData)
	if err != nil {
		return nil, fmt.Errorf("manifest.ListFiles: parse baseline: %w", err)
	}
	fileMap := make(map[string]FileEntry)
	for _, f := range baseDeck.F {
		fileMap[f.Name] = FileEntry{Name: f.Name, UUID: f.UUID, Perm: f.Perm}
	}
	for _, f := range d.F {
		if f.UUID == "" {
			delete(fileMap, f.Name)
		} else {
			fileMap[f.Name] = FileEntry{Name: f.Name, UUID: f.UUID, Perm: f.Perm}
		}
	}
	entries := make([]FileEntry, 0, len(fileMap))
	for _, e := range fileMap {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries, nil
}

func fCardsToEntries(fCards []deck.FileCard) []FileEntry {
	entries := make([]FileEntry, len(fCards))
	for i, f := range fCards {
		entries[i] = FileEntry{Name: f.Name, UUID: f.UUID, Perm: f.Perm}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries
}
