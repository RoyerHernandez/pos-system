<?php

class InventoryController{

	/*=============================================
	FIXED CATALOG - MOVEMENT MOTIVOS
	=============================================*/

	static public $entryMotivos = array("Compra", "Devolución de cliente", "Ajuste positivo");
	static public $exitMotivos = array("Merma", "Devolución a proveedor", "Ajuste negativo");
	static public $autoMotivos = array("Venta", "Cancelación venta");

	/*=============================================
	CREATE MOVEMENT (manual)
	=============================================*/

	static public function ctrCreateMovement($data){

		// Validate tipo
		if($data["tipo"] != "entrada" && $data["tipo"] != "salida"){
			return array("status" => "error", "message" => "Tipo de movimiento inválido");
		}

		// Validate motivo against catalog (exclude auto-generated)
		if($data["tipo"] == "entrada"){
			if(!in_array($data["motivo"], self::$entryMotivos)){
				return array("status" => "error", "message" => "Motivo inválido para entrada");
			}
		}else{
			if(!in_array($data["motivo"], self::$exitMotivos)){
				return array("status" => "error", "message" => "Motivo inválido para salida");
			}
		}

		// Validate cantidad
		if(!isset($data["cantidad"]) || intval($data["cantidad"]) <= 0){
			return array("status" => "error", "message" => "La cantidad debe ser mayor a 0");
		}

		// Validate product exists
		$product = ProductModel::mdlShowProducts("productos", "id", $data["id_producto"]);

		if(!$product){
			return array("status" => "error", "message" => "Producto no encontrado");
		}

		// Validate stock for exits
		if($data["tipo"] == "salida" && intval($product["stock"]) < intval($data["cantidad"])){
			return array("status" => "error", "message" => "Stock insuficiente. Stock actual: " . $product["stock"]);
		}

		// Prepare movement data
		$movementData = array(
			"id_producto" => $data["id_producto"],
			"id_usuario" => $data["id_usuario"],
			"tipo" => $data["tipo"],
			"motivo" => $data["motivo"],
			"cantidad" => intval($data["cantidad"]),
			"observaciones" => isset($data["observaciones"]) ? $data["observaciones"] : null,
			"id_referencia" => null
		);

		// Insert movement
		$response = InventoryModel::mdlInsertMovement("movimientos_inventario", $movementData);

		if($response != "ok"){
			return array("status" => "error", "message" => "Error al registrar el movimiento");
		}

		// Update stock
		$operation = ($data["tipo"] == "entrada") ? "add" : "subtract";
		ProductModel::mdlUpdateStock("productos", $data["id_producto"], intval($data["cantidad"]), $operation);

		return array("status" => "ok", "message" => "Movimiento registrado correctamente");

	}

	/*=============================================
	SHOW MOVEMENTS
	=============================================*/

	static public function ctrShowMovements($item, $value){

		$response = InventoryModel::mdlShowMovements("movimientos_inventario", $item, $value);

		return $response;

	}

	/*=============================================
	SHOW MOVEMENTS BY DATE RANGE
	=============================================*/

	static public function ctrShowMovementsByDate($startDate, $endDate){

		$response = InventoryModel::mdlShowMovementsByDate("movimientos_inventario", $startDate, $endDate);

		return $response;

	}

	/*=============================================
	GET MOTIVOS (for UI consumption)
	=============================================*/

	static public function ctrGetMotivos(){

		return array(
			"entrada" => self::$entryMotivos,
			"salida" => self::$exitMotivos
		);

	}

}
