package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/paramientos/leap/internal/config"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch web dashboard for managing connections",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		fmt.Println("\n🌐 \033[1;32mLEAP WEB DASHBOARD\033[0m")
		fmt.Println("\033[90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n")

		// Setup routes
		http.HandleFunc("/", serveHome)
		http.HandleFunc("/api/connections", apiConnections)
		http.HandleFunc("/api/stats", apiStats)

		addr := ":" + port
		fmt.Printf("  🚀 Server starting on \033[1;36mhttp://localhost%s\033[0m\n", addr)
		fmt.Println("\n\033[90m  Press Ctrl+C to stop the server\033[0m\n")

		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err)
		}
	},
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>⚡ LEAP Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        .header {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        .header h1 {
            font-size: 2.5em;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 10px;
        }
        .header p { color: #666; font-size: 1.1em; }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            transition: transform 0.3s ease;
        }
        .stat-card:hover { transform: translateY(-5px); }
        .stat-card h3 {
            color: #667eea;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 10px;
        }
        .stat-card .value { font-size: 2.5em; font-weight: bold; color: #333; }
        .connections {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 30px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        .connections h2 { color: #333; margin-bottom: 20px; font-size: 1.8em; }
        .connection-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 15px;
            color: white;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: transform 0.2s ease;
            cursor: pointer;
        }
        .connection-card:hover { transform: scale(1.02); }
        .connection-info h3 { font-size: 1.3em; margin-bottom: 5px; }
        .connection-info p { opacity: 0.9; font-size: 0.95em; }
        .connection-badge {
            background: rgba(255,255,255,0.2);
            padding: 8px 15px;
            border-radius: 20px;
            font-size: 0.85em;
            font-weight: 600;
        }
        .loading { text-align: center; padding: 40px; color: white; font-size: 1.2em; }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
        .pulse { animation: pulse 2s infinite; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚡ LEAP Dashboard</h1>
            <p>Manage your SSH connections from anywhere</p>
        </div>
        <div class="stats">
            <div class="stat-card">
                <h3>Total Connections</h3>
                <div class="value" id="total-connections">-</div>
            </div>
            <div class="stat-card">
                <h3>Favorites</h3>
                <div class="value" id="favorites">-</div>
            </div>
            <div class="stat-card">
                <h3>Last Used</h3>
                <div class="value" id="last-used">-</div>
            </div>
        </div>
        <div class="connections">
            <h2>📡 Your Connections</h2>
            <div id="connections-list" class="loading pulse">Loading connections...</div>
        </div>
    </div>
    <script>
        async function loadData() {
            try {
                const response = await fetch('/api/connections');
                const data = await response.json();
                document.getElementById('total-connections').textContent = data.total;
                document.getElementById('favorites').textContent = data.favorites;
                document.getElementById('last-used').textContent = data.last_used || 'Never';
                const list = document.getElementById('connections-list');
                if (data.connections.length === 0) {
                    list.innerHTML = '<p style="color: #666;">No connections found. Add one with <code>leap add</code></p>';
                } else {
                    list.innerHTML = data.connections.map(conn => ` + "`" + `
                        <div class="connection-card">
                            <div class="connection-info">
                                <h3>${conn.favorite ? '⭐ ' : ''}${conn.name}</h3>
                                <p>${conn.user}@${conn.host}:${conn.port}</p>
                            </div>
                            <div class="connection-badge">
                                ${conn.tags && conn.tags.length > 0 ? conn.tags.join(', ') : 'No tags'}
                            </div>
                        </div>
                    ` + "`" + `).join('');
                }
            } catch (error) {
                document.getElementById('connections-list').innerHTML = 
                    '<p style="color: #ff6b6b;">Error loading connections: ' + error.message + '</p>';
            }
        }
        loadData();
        setInterval(loadData, 5000);
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func apiConnections(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(os.Getenv("LEAP_MASTER_PASSWORD"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ConnectionResponse struct {
		Name     string   `json:"name"`
		Host     string   `json:"host"`
		User     string   `json:"user"`
		Port     int      `json:"port"`
		Tags     []string `json:"tags"`
		Favorite bool     `json:"favorite"`
		LastUsed string   `json:"last_used"`
	}

	var connections []ConnectionResponse
	var favoriteCount int
	var lastUsedName string
	var lastUsedTime time.Time

	for _, conn := range cfg.Connections {
		connections = append(connections, ConnectionResponse{
			Name:     conn.Name,
			Host:     conn.Host,
			User:     conn.User,
			Port:     conn.Port,
			Tags:     conn.Tags,
			Favorite: conn.Favorite,
			LastUsed: conn.LastUsed.Format("2006-01-02 15:04"),
		})

		if conn.Favorite {
			favoriteCount++
		}

		if conn.LastUsed.After(lastUsedTime) {
			lastUsedTime = conn.LastUsed
			lastUsedName = conn.Name
		}
	}

	response := map[string]interface{}{
		"connections": connections,
		"total":       len(connections),
		"favorites":   favoriteCount,
		"last_used":   lastUsedName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func apiStats(w http.ResponseWriter, r *http.Request) {
	// Placeholder for future stats endpoint
	stats := map[string]interface{}{
		"uptime": time.Since(time.Now()).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func init() {
	webCmd.Flags().StringP("port", "p", "8080", "Port to run web server on")
	rootCmd.AddCommand(webCmd)
}
