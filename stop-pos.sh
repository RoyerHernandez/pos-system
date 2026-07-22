#!/usr/bin/env bash
# Mata el servidor PHP de POS y limpia las sesiones

echo "Matando servidor PHP..."
pkill -f "php -S localhost:8000" 2>/dev/null && echo "Servidor detenido." || echo "No habia servidor corriendo."

echo "Limpiando sesiones PHP..."
find /tmp -maxdepth 1 -name "sess_*" -delete 2>/dev/null
echo "Sesiones eliminadas."

echo "Listo. Ahora reinicia con:"
echo "  cd pos && php -S localhost:8000 router.php"
