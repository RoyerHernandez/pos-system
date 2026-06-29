<?php

class SaleController{

	/*=============================================
	CREATE SALE (called via AJAX)
	=============================================*/

	static public function ctrCreateSale($saleData){

		// Generate sale code
		$saleCode = SaleModel::mdlGenerateSaleCode("ventas");

		// Get open cash register
		$cashRegister = CashRegisterModel::mdlGetOpenCashRegister("caja", $saleData["id_usuario"]);

		$idCaja = $cashRegister ? $cashRegister["id"] : null;

		// Insert sale
		$data = array(
			"id_usuario" => $saleData["id_usuario"],
			"id_cliente" => $saleData["id_cliente"],
			"id_caja" => $idCaja,
			"codigo_venta" => $saleCode,
			"subtotal" => $saleData["subtotal"],
			"impuesto" => $saleData["impuesto"],
			"descuento" => $saleData["descuento"],
			"total" => $saleData["total"],
			"metodo_pago" => $saleData["metodo_pago"]
		);

		$idSale = SaleModel::mdlInsertSale("ventas", $data);

		if($idSale == "error"){
			return array("status" => "error", "message" => "Error al registrar la venta");
		}

		// Insert sale details and update stock
		$products = json_decode($saleData["products"], true);

		foreach($products as $product){

			$detail = array(
				"id_venta" => $idSale,
				"id_producto" => $product["id"],
				"cantidad" => $product["quantity"],
				"precio_unitario" => $product["price"],
				"descuento" => $product["discount"],
				"subtotal" => $product["subtotal"]
			);

			$responseDetail = SaleModel::mdlInsertSaleDetail("detalle_ventas", $detail);

			if($responseDetail == "error"){
				return array("status" => "error", "message" => "Error al registrar detalle de venta");
			}

			// Update stock (subtract)
			ProductModel::mdlUpdateStock("productos", $product["id"], $product["quantity"], "subtract");

			// Register inventory exit movement
			InventoryModel::mdlInsertMovement("movimientos_inventario", array(
				"id_producto" => $product["id"],
				"id_usuario" => $saleData["id_usuario"],
				"tipo" => "salida",
				"motivo" => "Venta",
				"cantidad" => $product["quantity"],
				"observaciones" => "Venta " . $saleCode,
				"id_referencia" => $idSale
			));

		}

		// Update client purchases
		if($saleData["id_cliente"] > 0){
			ClientModel::mdlUpdatePurchases("clientes", $saleData["id_cliente"], $saleData["total"]);
		}

		// Update cash register totals
		if($idCaja){
			CashRegisterModel::mdlUpdateCashRegisterTotals("caja", $idCaja, $saleData["total"], $saleData["metodo_pago"]);
		}

		return array(
			"status" => "ok",
			"saleCode" => $saleCode,
			"saleId" => $idSale
		);

	}

	/*=============================================
	CANCEL SALE
	=============================================*/

	static public function ctrCancelSale($idSale){

		// Get sale details to restore stock
		$details = SaleModel::mdlShowSaleDetails("detalle_ventas", $idSale);

		foreach($details as $detail){
			// Restore stock (add back)
			ProductModel::mdlUpdateStock("productos", $detail["id_producto"], $detail["cantidad"], "add");

			// Register inventory entry movement (reversal)
			InventoryModel::mdlInsertMovement("movimientos_inventario", array(
				"id_producto" => $detail["id_producto"],
				"id_usuario" => $_SESSION["id"],
				"tipo" => "entrada",
				"motivo" => "Cancelación venta",
				"cantidad" => $detail["cantidad"],
				"observaciones" => "Cancelación venta #" . $idSale,
				"id_referencia" => $idSale
			));
		}

		// Get sale info to update client purchases
		$sale = SaleModel::mdlShowSales("ventas", "v.id", $idSale);

		if($sale && $sale["id_cliente"] > 0){
			ClientModel::mdlUpdatePurchases("clientes", $sale["id_cliente"], -$sale["total"]);
		}

		// Cancel the sale
		$response = SaleModel::mdlCancelSale("ventas", $idSale);

		return $response;

	}

	/*=============================================
	SHOW SALES
	=============================================*/

	static public function ctrShowSales($item, $value){

		$response = SaleModel::mdlShowSales("ventas", $item, $value);

		return $response;

	}

	/*=============================================
	SHOW SALE DETAILS
	=============================================*/

	static public function ctrShowSaleDetails($idSale){

		$response = SaleModel::mdlShowSaleDetails("detalle_ventas", $idSale);

		return $response;

	}

}
