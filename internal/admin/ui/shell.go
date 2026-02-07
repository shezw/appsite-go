package ui

const AdminShellHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Appsite Admin</title>
    
    <!-- ES Module Shim (for compatibility) -->
    <script async src="https://ga.jspm.io/npm:es-module-shims@1.7.0/dist/es-module-shims.js"></script>
    
    <!-- Import Map to manage external dependencies neatly if needed, or straight imports in JS -->
    <script type="importmap">
    {
        "imports": {
            "react": "https://esm.sh/react@18.2.0",
            "react-dom/client": "https://esm.sh/react-dom@18.2.0/client"
        }
    }
    </script>

    <style>body { margin: 0; padding: 0; font-family: sans-serif; }</style>
</head>
<body>
    <div id="root"></div>

    <!-- The Application Logic -->
    <!-- We serve src/App.jsx as a static file, but browsers need Babel if we use JSX. 
         Ideally, we should pre-compile.
         However, since we cannot run node in this environment, we use Babel Standalone to transform on the fly.
    -->
    <script src="https://unpkg.com/@babel/standalone/babel.min.js"></script>
    
    <script type="text/babel" data-type="module" data-presets="react" src="/admin-assets/src/App.jsx"></script>
</body>
</html>`
