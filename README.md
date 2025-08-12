# 🤪 Generador de Excusas Divertidas

Un proyecto simple pero divertido que demuestra la creación de una aplicación web completa, desde la **API backend** con Go y el framework Gin hasta una **interfaz de usuario (UI)** con HTML, CSS y JavaScript.

## ✨ Características Principales

* **Generador de Excusas:** Obtén una excusa aleatoria y divertida con solo un clic, categorizada por el evento que desees.
* **API RESTful:** Endpoints sencillos para generar (`/api/generate`) y buscar (`/api/search`) excusas.
* **Persistencia de Datos:** Utiliza una base de datos **SQLite** para almacenar las excusas generadas, con ejemplos iniciales precargados.
* **Optimización de Rendimiento:** Implementa **caching** en memoria (`go-cache`) para búsquedas frecuentes y un **rate limiter** simple para proteger la API.
* **Frontend Interactivo:** Una UI minimalista y responsive, construida con HTML y el framework CSS Pico.css, que permite interactuar con la API de forma dinámica.

## 🚀 Cómo Empezar

Este proyecto está diseñado para ser fácil de configurar y ejecutar. Asegúrate de tener **Go** instalado.

1.  **Clona el repositorio:**
    ```bash
    git clone https://github.com/santiagourdaneta/Generador-de-Excusas-Divertidas-/
    cd Generador-de-Excusas-Divertidas-
    ```

2.  **Instala las dependencias:**
    ```bash
    go mod tidy
    ```
    Este comando descargará `gin-gonic/gin`, `modernc.org/sqlite` y `patrickmn/go-cache`.

3.  **Ejecuta la aplicación:**
    ```bash
    go run main.go
    ```
    La aplicación se ejecutará en `http://localhost:8080`. Abre tu navegador y ¡a generar excusas!

## 🛠️ Tecnologías Usadas

* **Backend:**
    * **Go:** El lenguaje de programación principal.
    * **Gin Gonic:** Un framework web de alto rendimiento para Go.
    * **SQLite:** Una base de datos ligera y autónoma para la persistencia de datos.
    * **go-cache:** Una implementación de caché en memoria para Go.
* **Frontend:**
    * **HTML, CSS, JavaScript:** Los bloques de construcción de la web.
    * **Pico.css:** Un framework CSS minimalista y sin clases para un estilo rápido y limpio.

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Si tienes ideas para nuevas excusas, mejoras en el código o cualquier otra cosa, no dudes en abrir un *issue* o enviar un *pull request*.

---

🏷️ Etiquetas y Hashtags
Labels de GitHub: go, golang, gin, sqlite, web-api, api, fullstack, web-app, html, css, javascript, fun

Hashtags: #GoLang #GinFramework #WebDev #API #FullStack #SQLite
