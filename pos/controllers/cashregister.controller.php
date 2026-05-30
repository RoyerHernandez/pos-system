<?php

class CashRegisterController{

	/*=============================================
	OPEN CASH REGISTER
	=============================================*/

	static public function ctrOpenCashRegister(){

		if(isset($_POST["openingAmount"])){

			$data = array(
				"id_usuario" => $_SESSION["id"],
				"monto_apertura" => $_POST["openingAmount"]
			);

			$response = CashRegisterModel::mdlOpenCashRegister("caja", $data);

			return $response;

		}

	}

	/*=============================================
	CLOSE CASH REGISTER
	=============================================*/

	static public function ctrCloseCashRegister(){

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

				return $response;

			}

		}

	}

	/*=============================================
	GET OPEN CASH REGISTER
	=============================================*/

	static public function ctrGetOpenCashRegister(){

		$response = CashRegisterModel::mdlGetOpenCashRegister("caja", $_SESSION["id"]);

		return $response;

	}

}
