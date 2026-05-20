package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"yigou-ide/backend/models"

	_ "modernc.org/sqlite"
)

type ProjectService struct {
	mu       sync.RWMutex
	db       *sql.DB
	dbPath   string
	projects map[int64]*models.Project
}

func NewProjectService() *ProjectService {
	homeDir, _ := os.UserHomeDir()
	dbDir := filepath.Join(homeDir, ".yigou-ide")
	os.MkdirAll(dbDir, 0755)
	dbPath := filepath.Join(dbDir, "projects.db")

	ps := &ProjectService{
		dbPath:   dbPath,
		projects: make(map[int64]*models.Project),
	}

	ps.initDB()
	ps.loadProjects()
	return ps
}

func (ps *ProjectService) initDB() {
	var err error
	ps.db, err = sql.Open("sqlite", ps.dbPath)
	if err != nil {
		fmt.Printf("打开数据库失败: %v\n", err)
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		path TEXT NOT NULL UNIQUE,
		type TEXT DEFAULT 'window',
		author TEXT DEFAULT '',
		version TEXT DEFAULT '1.0.0',
		icon TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_opened DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_favorite INTEGER DEFAULT 0,
		tags TEXT DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS project_templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		category TEXT DEFAULT '通用',
		icon TEXT DEFAULT '',
		content TEXT DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS recent_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		file_path TEXT NOT NULL,
		opened_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);
	`

	_, err = ps.db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("创建数据库表失败: %v\n", err)
		return
	}

	ps.initDefaultTemplates()
}

func (ps *ProjectService) initDefaultTemplates() {
	var count int
	ps.db.QueryRow("SELECT COUNT(*) FROM project_templates").Scan(&count)
	if count > 0 {
		return
	}

	templates := []struct {
		name        string
		description string
		category    string
		icon        string
		content     string
	}{
		{
			name:        "空白窗口程序",
			description: "创建一个空白的窗口应用程序",
			category:    "基础",
			icon:        "window",
			content:     `{"appName":"新项目","windowTitle":"新窗口","width":800,"height":600,"resizable":true,"minWidth":400,"minHeight":300}`,
		},
		{
			name:        "Hello World",
			description: "经典的Hello World示例程序",
			category:    "示例",
			icon:        "hello",
			content:     `{"appName":"HelloWorld","windowTitle":"Hello World","width":600,"height":400,"resizable":true}`,
		},
		{
			name:        "计算器",
			description: "简单的计算器应用程序模板",
			category:    "工具",
			icon:        "calculator",
			content:     `{"appName":"计算器","windowTitle":"计算器","width":350,"height":500,"resizable":false,"minWidth":350,"minHeight":500}`,
		},
		{
			name:        "文本编辑器",
			description: "简易文本编辑器模板",
			category:    "工具",
			icon:        "editor",
			content:     `{"appName":"文本编辑器","windowTitle":"文本编辑器","width":800,"height":600,"resizable":true,"minWidth":600,"minHeight":400}`,
		},
		{
			name:        "数据管理",
			description: "数据增删改查管理模板",
			category:    "数据",
			icon:        "database",
			content:     `{"appName":"数据管理","windowTitle":"数据管理系统","width":1000,"height":700,"resizable":true,"minWidth":800,"minHeight":600}`,
		},
		{
			name:        "网络请求工具",
			description: "HTTP请求测试工具模板",
			category:    "网络",
			icon:        "network",
			content:     `{"appName":"网络工具","windowTitle":"网络请求工具","width":700,"height":500,"resizable":true}`,
		},
	}

	for _, t := range templates {
		ps.db.Exec(
			"INSERT INTO project_templates (name, description, category, icon, content) VALUES (?, ?, ?, ?, ?)",
			t.name, t.description, t.category, t.icon, t.content,
		)
	}
}

func (ps *ProjectService) loadProjects() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.db == nil {
		return
	}

	rows, err := ps.db.Query("SELECT id, name, description, path, type, author, version, icon, created_at, updated_at, last_opened, is_favorite, tags FROM projects ORDER BY last_opened DESC")
	if err != nil {
		return
	}
	defer rows.Close()

	ps.projects = make(map[int64]*models.Project)
	for rows.Next() {
		p := &models.Project{}
		var isFavorite int
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Path, &p.Type, &p.Author, &p.Version, &p.Icon, &p.CreatedAt, &p.UpdatedAt, &p.LastOpened, &isFavorite, &p.Tags)
		if err != nil {
			continue
		}
		p.IsFavorite = isFavorite == 1
		ps.projects[p.ID] = p
	}
}

func (ps *ProjectService) CreateProject(name, description, projectType, author, path string) (*models.Project, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	now := time.Now()
	result, err := ps.db.Exec(
		"INSERT INTO projects (name, description, path, type, author, version, created_at, updated_at, last_opened) VALUES (?, ?, ?, ?, ?, '1.0.0', ?, ?, ?)",
		name, description, path, projectType, author, now, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建项目失败: %w", err)
	}

	id, _ := result.LastInsertId()
	project := &models.Project{
		ID:          id,
		Name:        name,
		Description: description,
		Path:        path,
		Type:        projectType,
		Author:      author,
		Version:     "1.0.0",
		CreatedAt:   now,
		UpdatedAt:   now,
		LastOpened:  now,
	}

	ps.projects[id] = project

	projectDir := filepath.Join(path, name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("创建项目目录失败: %w", err)
	}

	srcDir := filepath.Join(projectDir, "src")
	buildDir := filepath.Join(projectDir, "build")
	intDir := filepath.Join(projectDir, "_int")
	for _, d := range []string{srcDir, buildDir, intDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, fmt.Errorf("创建项目子目录失败: %w", err)
		}
	}

	egoFile := filepath.Join(srcDir, name+".ego")
	egoData := map[string]interface{}{
		"version":    "1.0.0",
		"appName":    name,
		"author":     author,
		"type":       projectType,
		"components": []interface{}{},
		"codeRows": []map[string]string{
			{"event": "窗口.加载完成", "code": ""},
		},
	}
	data, _ := json.MarshalIndent(egoData, "", "  ")
	os.WriteFile(egoFile, data, 0644)

	return project, nil
}

func (ps *ProjectService) OpenProject(id int64) (*models.Project, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return nil, fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET last_opened = ? WHERE id = ?", now, id)
	project.LastOpened = now

	return project, nil
}

func (ps *ProjectService) GetProject(id int64) (*models.Project, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	project, exists := ps.projects[id]
	if !exists {
		return nil, fmt.Errorf("项目不存在: %d", id)
	}
	return project, nil
}

func (ps *ProjectService) SetProjectName(id int64, name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET name = ?, updated_at = ? WHERE id = ?", name, now, id)
	project.Name = name
	project.UpdatedAt = now
	return nil
}

func (ps *ProjectService) SetProjectDescription(id int64, description string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET description = ?, updated_at = ? WHERE id = ?", description, now, id)
	project.Description = description
	project.UpdatedAt = now
	return nil
}

func (ps *ProjectService) SetProjectFavorite(id int64, favorite bool) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	project.IsFavorite = favorite
	favInt := 0
	if favorite {
		favInt = 1
	}
	ps.db.Exec("UPDATE projects SET is_favorite = ? WHERE id = ?", favInt, id)
	return nil
}

func (ps *ProjectService) SetProjectTags(id int64, tags string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET tags = ?, updated_at = ? WHERE id = ?", tags, now, id)
	project.Tags = tags
	project.UpdatedAt = now
	return nil
}

func (ps *ProjectService) SetProjectVersion(id int64, version string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET version = ?, updated_at = ? WHERE id = ?", version, now, id)
	project.Version = version
	project.UpdatedAt = now
	return nil
}

func (ps *ProjectService) GetProjectConfig(id int64) (*models.ProjectConfig, error) {
	ps.mu.RLock()
	project, exists := ps.projects[id]
	if !exists {
		ps.mu.RUnlock()
		return nil, fmt.Errorf("项目不存在: %d", id)
	}
	projectPath := project.Path
	projectName := project.Name
	ps.mu.RUnlock()

	egoFile := filepath.Join(projectPath, projectName, "src", projectName+".ego")
	if _, err := os.Stat(egoFile); os.IsNotExist(err) {
		egoFile = filepath.Join(projectPath, projectName, projectName+".ego")
	}
	data, err := os.ReadFile(egoFile)
	if err != nil {
		return nil, fmt.Errorf("读取项目配置失败: %w", err)
	}

	var config models.ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析项目配置失败: %w", err)
	}
	return &config, nil
}

func (ps *ProjectService) SetProjectConfig(id int64, config *models.ProjectConfig) error {
	ps.mu.RLock()
	project, exists := ps.projects[id]
	if !exists {
		ps.mu.RUnlock()
		return fmt.Errorf("项目不存在: %d", id)
	}
	projectPath := project.Path
	projectName := project.Name
	ps.mu.RUnlock()

	srcDir := filepath.Join(projectPath, projectName, "src")
	egoFile := filepath.Join(srcDir, projectName+".ego")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		egoFile = filepath.Join(projectPath, projectName, projectName+".ego")
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化项目配置失败: %w", err)
	}

	if err := os.WriteFile(egoFile, data, 0644); err != nil {
		return fmt.Errorf("写入项目配置失败: %w", err)
	}

	now := time.Now()
	ps.mu.Lock()
	ps.db.Exec("UPDATE projects SET updated_at = ? WHERE id = ?", now, id)
	if p, ok := ps.projects[id]; ok {
		p.UpdatedAt = now
	}
	ps.mu.Unlock()

	return nil
}

func (ps *ProjectService) DeleteProject(id int64) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, err := ps.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	delete(ps.projects, id)
	return nil
}

func (ps *ProjectService) GetAllProjects() []*models.Project {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	projects := make([]*models.Project, 0, len(ps.projects))
	for _, p := range ps.projects {
		projects = append(projects, p)
	}
	return projects
}

func (ps *ProjectService) GetRecentProjects(limit int) []*models.Project {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	projects := make([]*models.Project, 0)
	for _, p := range ps.projects {
		projects = append(projects, p)
		if len(projects) >= limit {
			break
		}
	}
	return projects
}

func (ps *ProjectService) GetFavoriteProjects() []*models.Project {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	projects := make([]*models.Project, 0)
	for _, p := range ps.projects {
		if p.IsFavorite {
			projects = append(projects, p)
		}
	}
	return projects
}

func (ps *ProjectService) ToggleFavorite(id int64) (bool, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return false, fmt.Errorf("项目不存在: %d", id)
	}

	project.IsFavorite = !project.IsFavorite
	favInt := 0
	if project.IsFavorite {
		favInt = 1
	}

	ps.db.Exec("UPDATE projects SET is_favorite = ? WHERE id = ?", favInt, id)
	return project.IsFavorite, nil
}

func (ps *ProjectService) UpdateProject(id int64, name, description string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	project, exists := ps.projects[id]
	if !exists {
		return fmt.Errorf("项目不存在: %d", id)
	}

	now := time.Now()
	ps.db.Exec("UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?", name, description, now, id)

	project.Name = name
	project.Description = description
	project.UpdatedAt = now

	return nil
}

func (ps *ProjectService) GetTemplates() []models.ProjectTemplate {
	rows, err := ps.db.Query("SELECT id, name, description, category, icon, content FROM project_templates ORDER BY category, name")
	if err != nil {
		return nil
	}
	defer rows.Close()

	templates := make([]models.ProjectTemplate, 0)
	for rows.Next() {
		t := models.ProjectTemplate{}
		rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category, &t.Icon, &t.Content)
		templates = append(templates, t)
	}
	return templates
}

func (ps *ProjectService) SearchProjects(query string) []*models.Project {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	projects := make([]*models.Project, 0)
	for _, p := range ps.projects {
		if contains(p.Name, query) || contains(p.Description, query) || contains(p.Tags, query) {
			projects = append(projects, p)
		}
	}
	return projects
}

func (ps *ProjectService) Close() {
	if ps.db != nil {
		ps.db.Close()
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
