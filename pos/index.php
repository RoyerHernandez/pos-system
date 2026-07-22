<?php

require_once __DIR__ . "/controllers/template.controller.php";
require_once __DIR__ . "/controllers/users.controller.php";
require_once __DIR__ . "/controllers/categories.controller.php";
require_once __DIR__ . "/controllers/products.controller.php";
require_once __DIR__ . "/controllers/clients.controller.php";
require_once __DIR__ . "/controllers/sales.controller.php";
require_once __DIR__ . "/controllers/cashregister.controller.php";
require_once __DIR__ . "/controllers/dashboard.controller.php";
require_once __DIR__ . "/controllers/reports.controller.php";
require_once __DIR__ . "/controllers/inventory.controller.php";

require_once __DIR__ . "/models/users.model.php";
require_once __DIR__ . "/models/categories.model.php";
require_once __DIR__ . "/models/products.model.php";
require_once __DIR__ . "/models/clients.model.php";
require_once __DIR__ . "/models/sales.model.php";
require_once __DIR__ . "/models/cashregister.model.php";
require_once __DIR__ . "/models/dashboard.model.php";
require_once __DIR__ . "/models/reports.model.php";
require_once __DIR__ . "/models/inventory.model.php";

$template = new TemplateController();
$template -> ctrTemplate();