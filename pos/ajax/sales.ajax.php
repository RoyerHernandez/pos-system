<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("status" => "error", "message" => "No autorizado"));
	return;
}

require_once "../models/connection.php";
require_once "../models/sales.model.php";
require_once "../models/products.model.php";
require_once "../models/clients.model.php";
require_once "../models/cashregister.model.php";

require_once "../controllers/sales.controller.php";
require_once "../controllers/cashregister.controller.php";

/*=============================================
COMPLETE SALE
=============================================*/

if(isset($_POST["completeSale"])){

	$saleData = array(
		"id_usuario" => $_SESSION["id"],
		"id_cliente" => $_POST["idClient"],
		"products" => $_POST["products"],
		"subtotal" => $_POST["subtotal"],
		"impuesto" => $_POST["tax"],
		"descuento" => $_POST["discount"],
		"total" => $_POST["total"],
		"metodo_pago" => $_POST["paymentMethod"]
	);

	$response = SaleController::ctrCreateSale($saleData);

	echo json_encode($response);

}

/*=============================================
CANCEL SALE
=============================================*/

if(isset($_POST["cancelSale"])){

	$response = SaleController::ctrCancelSale($_POST["idSale"]);

	echo json_encode(array("status" => $response));

}

/*=============================================
GET SALE DETAILS
=============================================*/

if(isset($_POST["getSaleDetails"])){

	$details = SaleController::ctrShowSaleDetails($_POST["idSale"]);

	echo json_encode($details);

}

/*=============================================
GET ALL SALES (for DataTable)
=============================================*/

if(isset($_POST["getSales"])){

	if(isset($_POST["startDate"]) && isset($_POST["endDate"]) && $_POST["startDate"] != "" && $_POST["endDate"] != ""){

		$response = SaleModel::mdlShowSalesByDate("ventas", $_POST["startDate"], $_POST["endDate"]);

	}else{

		$response = SaleController::ctrShowSales(null, null);

	}

	echo json_encode($response);

}

/*=============================================
OPEN CASH REGISTER
=============================================*/

if(isset($_POST["openCashRegister"])){

	$data = array(
		"id_usuario" => $_SESSION["id"],
		"monto_apertura" => $_POST["openingAmount"]
	);

	$response = CashRegisterModel::mdlOpenCashRegister("caja", $data);

	echo json_encode(array("status" => $response));

}

/*=============================================
CLOSE CASH REGISTER
=============================================*/

if(isset($_POST["closeCashRegister"])){

	$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);

	if($cashRegister){

		$data = array(
			"id" => $cashRegister["id"],
			"monto_cierre" => $_POST["closingAmount"],
			"total_ventas" => $cashRegister["total_ventas"],
			"total_efectivo" => $cashRegister["total_efectivo"],
			"total_tarjeta" => $cashRegister["total_tarjeta"],
			"total_transferencia" => $cashRegister["total_transferencia"]
		);

		$response = CashRegisterModel::mdlCloseCashRegister("caja", $data);

		echo json_encode(array("status" => $response, "cashRegister" => $cashRegister));

	}else{

		echo json_encode(array("status" => "error", "message" => "No hay caja abierta"));

	}

}

/*=============================================
GET CASH REGISTER STATUS
=============================================*/

if(isset($_POST["getCashRegisterStatus"])){

	$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);

	if($cashRegister){
		echo json_encode(array("status" => "open", "cashRegister" => $cashRegister));
	}else{
		echo json_encode(array("status" => "closed"));
	}

}
