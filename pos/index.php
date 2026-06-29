<?php

require_once "controllers/template.controller.php";
require_once "controllers/users.controller.php";
require_once "controllers/categories.controller.php";
require_once "controllers/products.controller.php";
require_once "controllers/clients.controller.php";
require_once "controllers/sales.controller.php";
require_once "controllers/cashregister.controller.php";
require_once "controllers/dashboard.controller.php";
require_once "controllers/reports.controller.php";

require_once "models/users.model.php";
require_once "models/categories.model.php";
require_once "models/products.model.php";
require_once "models/clients.model.php";
require_once "models/sales.model.php";
require_once "models/cashregister.model.php";
require_once "models/dashboard.model.php";
require_once "models/reports.model.php";

$template = new TemplateController();
$template -> ctrTemplate();