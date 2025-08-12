package main

import (
    "database/sql"
    "fmt"
    "math/rand"
    "net/http"
    "strings"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
    _ "modernc.org/sqlite"
    "github.com/patrickmn/go-cache"
    "html"
)

var db *sql.DB
var c *cache.Cache

func main() {
    var err error
    db, err = sql.Open("sqlite", "./excuses.db")
    if err != nil { panic(err) }
    defer db.Close()

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS excuses (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            excuse_text TEXT NOT NULL CHECK (LENGTH(excuse_text) <= 200),
            category TEXT CHECK (LENGTH(category) <= 50),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_text ON excuses(excuse_text);
    `)
    if err != nil { panic("Can't create table/index: " + err.Error()) }

    _, err = db.Exec(`
        INSERT OR IGNORE INTO excuses (excuse_text, category) VALUES
        ('My cat ate my homework', 'school'),
        ('My dog broke my phone', 'party'),
        ('My bird hid my keys', 'dinner');
    `)
    if err != nil { panic("Can't insert sample data: " + err.Error()) }

    err = db.Ping()
    if err != nil { panic("Can't connect to database: " + err.Error()) }

    c = cache.New(5*time.Minute, 10*time.Minute)

    rows, err := db.Query("SELECT excuse_text FROM excuses WHERE excuse_text LIKE ?", "%party%")
    if err == nil {
        var results []string
        for rows.Next() {
            var text string
            rows.Scan(&text)
            results = append(results, text)
        }
        rows.Close()
        c.Set("search_party", results, cache.DefaultExpiration)
    }

    r := gin.Default()
    r.Use(gin.Recovery())
    r.Use(rateLimiterMiddleware())
    r.Static("/static", "./static")
    r.GET("/api/generate", generateExcuse)
    r.GET("/api/search", searchExcuses)
    r.GET("/", homePage)

    r.Run(":8080")
}

func generateExcuse(ctx *gin.Context) {
    category := ctx.Query("category")
    if len(category) > 50 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category too long (max 50 chars)"})
        return
    }
    if category == "" {
        category = "party"
    }
    category = strings.TrimSpace(html.EscapeString(category))

    animals := []string{"cat", "dog", "bird", "fish"}
    actions := []string{"ate", "hid", "broke", "stole"}
    things := []string{"shoes", "phone", "car keys", "homework"}
    rand.Seed(time.Now().UnixNano())
    excuse := fmt.Sprintf("My %s %s my %s!", animals[rand.Intn(len(animals))], actions[rand.Intn(len(actions))], things[rand.Intn(len(things))])

    stmt, err := db.Prepare("INSERT INTO excuses (excuse_text, category) VALUES (?, ?)")
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }
    defer stmt.Close()
    _, err = stmt.Exec(excuse, category)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save excuse: " + err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"excuse": excuse})
}

func searchExcuses(ctx *gin.Context) {
    query := ctx.Query("q")
    if query == "" { ctx.JSON(http.StatusBadRequest, gin.H{"error": "No query"}); return }

    query = strings.TrimSpace(html.EscapeString(query))

    pageStr := ctx.Query("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    offset := (page - 1) * 10

    cacheKey := "search_" + query + "_page_" + fmt.Sprint(page)
    if val, found := c.Get(cacheKey); found {
        ctx.JSON(http.StatusOK, val)
        return
    }

    rows, err := db.Query("SELECT excuse_text FROM excuses WHERE excuse_text LIKE ? LIMIT 10 OFFSET ?", "%"+query+"%", offset)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Search error: " + err.Error()})
        return
    }
    defer rows.Close()

    var results []string
    for rows.Next() {
        var text string
        rows.Scan(&text)
        results = append(results, text)
    }

    c.Set(cacheKey, results, cache.DefaultExpiration)

    ctx.JSON(http.StatusOK, results)
}

func rateLimiterMiddleware() gin.HandlerFunc {
    limiter := make(map[string]time.Time)
    return func(c *gin.Context) {
        ip := c.ClientIP()
        if t, ok := limiter[ip]; ok && time.Since(t) < time.Second {
            c.AbortWithStatus(http.StatusTooManyRequests)
            return
        }
        limiter[ip] = time.Now()
        c.Next()
    }
}

func homePage(ctx *gin.Context) {
    html := `
    <!DOCTYPE html>
    <html lang="es" data-theme="dark">
    <head>
        <meta charset="UTF-8">
        <title>Generador de Excusas</title>
        <link rel="stylesheet" href="/static/pico.css">
        <meta property="og:title" content="Excuse Maker">
        <meta property="og:description" content="Create fun excuses to skip events!">
        <meta property="og:image" content="https://example.com/excuse-image.jpg">
        <meta property="og:url" content="https://your-site.com">
        <meta name="twitter:card" content="summary_large_image">
        <meta name="description" content="Generate funny excuses instantly!">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>
            body { background-color: #1a1a1a; font-family: Arial, sans-serif; color: #e0e0e0; }
            .container { max-width: 600px; margin: 0 auto; padding: 20px; }
            h1 { color: #4dabf7; text-align: center; font-size: 2.5em; }
            input, button { font-size: 1.2em; padding: 12px; margin: 10px 0; border-radius: 5px; }
            button { background-color: #4caf50; color: white; border: none; cursor: pointer; width: 100%; }
            button:hover { background-color: #45a049; }
            #results { background-color: #2c2c2c; padding: 15px; border-radius: 5px; min-height: 50px; color: #e0e0e0; }
            #list li { font-size: 1.1em; padding: 8px; background-color: #333; margin: 5px 0; border-radius: 3px; }
            #load-more { display: none; margin: 15px auto; background-color: #0288d1; }
            #load-more:hover { background-color: #0277bd; }
            .error { color: #ff6b6b; font-size: 1.1em; }
            label { color: #b0bec5; }
        </style>
    </head>
    <body>
        <main class="container">
            <h1>🌙 Generador de Excusas 🌙</h1>
            <p style="text-align: center; color: #b0bec5;">¡Crea excusas divertidas para cualquier evento!</p>
            <form id="excuse-form" aria-label="Generar una excusa">
                <label for="category" style="font-size: 1.2em;">Evento (máx. 50 caracteres):</label>
                <input type="text" id="category" name="category" maxlength="50" placeholder="Ejemplo: fiesta" required aria-required="true" style="width: 100%; background-color: #333; color: #e0e0e0; border: 1px solid #555;">
                <button type="submit" aria-label="Generar excusa">¡Generar Excusa!</button>
            </form>
            <div id="results" aria-live="polite"></div>
            <label for="search" style="font-size: 1.2em;">Busca excusas (máx. 200 caracteres):</label>
            <input type="search" id="search" maxlength="200" placeholder="Ejemplo: gato" aria-label="Buscar excusas" style="width: 100%; background-color: #333; color: #e0e0e0; border: 1px solid #555;">
            <ul id="list" aria-label="Resultados de búsqueda"></ul>
            <button id="load-more" aria-label="Cargar más resultados">Cargar Más</button>
        </main>
        <script>
            const form = document.getElementById('excuse-form');
            const results = document.getElementById('results');
            const searchInput = document.getElementById('search');
            const list = document.getElementById('list');
            const loadMore = document.getElementById('load-more');
            let page = 1;

            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const category = document.getElementById('category').value;
                if (category.length > 50) {
                    results.innerHTML = '<p class="error">¡El evento es demasiado largo!</p>';
                    return;
                }
                results.innerHTML = '<p>Generando...</p>';
                try {
                    const res = await fetch('/api/generate?category=' + encodeURIComponent(category));
                    const data = await res.json();
                    if (data.error) {
                        results.innerHTML = '<p class="error">' + data.error + '</p>';
                    } else {
                        results.innerHTML = '<p><strong>¡Tu excusa!</strong> ' + data.excuse + '</p>';
                    }
                } catch (err) {
                    results.innerHTML = '<p class="error">¡Ups! Algo salió mal.</p>';
                }
            });

            function debounce(func, wait) {
                let timeout;
                return function executedFunction(...args) {
                    const later = () => {
                        clearTimeout(timeout);
                        func(...args);
                    };
                    clearTimeout(timeout);
                    timeout = setTimeout(later, wait);
                };
            }

            searchInput.addEventListener('input', debounce(async () => {
                const q = searchInput.value;
                if (q.length < 3) {
                    list.innerHTML = '';
                    loadMore.style.display = 'none';
                    return;
                }
                list.innerHTML = '<li>Buscando...</li>';
                try {
                    const res = await fetch('/api/search?q=' + encodeURIComponent(q) + '&page=' + page);
                    const data = await res.json();
                    list.innerHTML = '';
                    if (data.length === 0) {
                        list.innerHTML = '<li>No se encontraron excusas.</li>';
                    } else {
                        data.forEach(item => {
                            const li = document.createElement('li');
                            li.textContent = item;
                            list.appendChild(li);
                        });
                        loadMore.style.display = data.length === 10 ? 'block' : 'none';
                    }
                } catch (err) {
                    list.innerHTML = '<li class="error">¡Error al buscar!</li>';
                }
            }, 300));

            loadMore.addEventListener('click', async () => {
                page++;
                const q = searchInput.value;
                if (q.length < 3) return;
                try {
                    const res = await fetch('/api/search?q=' + encodeURIComponent(q) + '&page=' + page);
                    const data = await res.json();
                    if (data.length === 0) {
                        loadMore.style.display = 'none';
                    } else {
                        data.forEach(item => {
                            const li = document.createElement('li');
                            li.textContent = item;
                            list.appendChild(li);
                        });
                        loadMore.style.display = data.length === 10 ? 'block' : 'none';
                    }
                } catch (err) {
                    list.innerHTML += '<li class="error">¡Error al cargar más!</li>';
                }
            });
        </script>
    </body>
    </html>`
    ctx.Header("Content-Type", "text/html")
    ctx.String(http.StatusOK, html)
}