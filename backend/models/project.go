package models

import "time"

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	LastOpened  time.Time `json:"lastOpened"`
	IsFavorite  bool      `json:"isFavorite"`
	Tags        string    `json:"tags"`
}

type ProjectTemplate struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	Content     string `json:"content"`
}

type ProjectConfig struct {
	AppName     string `json:"appName"`
	WindowTitle string `json:"windowTitle"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Resizable   bool   `json:"resizable"`
	MinWidth    int    `json:"minWidth"`
	MinHeight   int    `json:"minHeight"`
	MaxWidth    int    `json:"maxWidth"`
	MaxHeight   int    `json:"maxHeight"`
	Fullscreen  bool   `json:"fullscreen"`
	Frameless   bool   `json:"frameless"`
	AlwaysOnTop bool   `json:"alwaysOnTop"`
	Background  string `json:"background"`
	Icon        string `json:"icon"`
}
