<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("error" => "No autorizado"));
	return;
}

// Only Admin and Especial can access reports
if(!isset($_SESSION["role"]) || ($_SESSION["role"] != "Administrador" && $_SESSION["role"] != "Especial")){
	echo json_encode(array("error" => "Acceso denegado"));
	return;
}

require_once "../models/connection.php";
require_once "../models/reports.model.php";
require_once "../controllers/reports.controller.php";

/*=============================================
SALES BY DATE RANGE
=============================================*/

if(isset($_POST["reportSalesByDate"])){

	$startDate = $_POST["startDate"];
	$endDate = $_POST["endDate"];

	$sales = ReportController::ctrSalesByDateRange($startDate, $endDate);
	$summary = ReportController::ctrSalesSummary($startDate, $endDate);
	$daily = ReportController::ctrDailySalesBreakdown($startDate, $endDate);

	echo json_encode(array(
		"sales" => $sales,
		"summary" => $summary,
		"daily" => $daily
	));

}

/*=============================================
TOP PRODUCTS
=============================================*/

if(isset($_POST["reportTopProducts"])){

	$startDate = $_POST["startDate"];
	$endDate = $_POST["endDate"];
	$limit = isset($_POST["limit"]) ? intval($_POST["limit"]) : 10;

	$products = ReportController::ctrTopProducts($startDate, $endDate, $limit);

	echo json_encode($products);

}

/*=============================================
SALES BY PAYMENT METHOD
=============================================*/

if(isset($_POST["reportByPaymentMethod"])){

	$startDate = $_POST["startDate"];
	$endDate = $_POST["endDate"];

	$data = ReportController::ctrSalesByPaymentMethod($startDate, $endDate);

	echo json_encode($data);

}

/*=============================================
CASH REGISTER REPORT
=============================================*/

if(isset($_POST["reportCashRegister"])){

	$startDate = $_POST["startDate"];
	$endDate = $_POST["endDate"];

	$data = ReportController::ctrCashRegisterReport($startDate, $endDate);

	echo json_encode($data);

}
