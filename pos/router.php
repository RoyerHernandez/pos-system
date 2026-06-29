<?php
/*=============================================
ROUTER FOR PHP BUILT-IN SERVER
Usage: php -S localhost:8080 router.php
Replaces .htaccess mod_rewrite rules
=============================================*/

$uri = parse_url($_SERVER["REQUEST_URI"], PHP_URL_PATH);

// Serve static files directly (css, js, images, fonts, etc.)
if($uri !== "/" && file_exists(__DIR__ . $uri)){
    return false;
}

// Route clean URLs: /productos -> index.php?ruta=productos
if(preg_match('#^/([-a-zA-Z0-9]+)$#', $uri, $matches)){
    $_GET["ruta"] = $matches[1];
}

require __DIR__ . "/index.php";
