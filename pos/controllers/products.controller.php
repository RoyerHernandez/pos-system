<?php

class ProductController{

	/*=============================================
	SHOW PRODUCTS
	=============================================*/

	static public function ctrShowProducts($item, $value){

		$table = "productos";

		$response = ProductModel::mdlShowProducts($table, $item, $value);

		return $response;

	}

	/*=============================================
	CREATE PRODUCT
	=============================================*/

	static public function ctrCreateProduct(){

		if(isset($_POST["nuevoCodigo"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevaDescripcion"])){

				/*=============================================
				VALIDATE IMAGE
				=============================================*/

				$path = "";

				if(isset($_FILES["nuevaImagen"]["tmp_name"]) && !empty($_FILES["nuevaImagen"]["tmp_name"])){

					list($width, $height) = getimagesize($_FILES["nuevaImagen"]["tmp_name"]);

					$newWidth = 500;
					$newHeight = 500;

					$absDir = __DIR__ . "/../views/img/productos/" . $_POST["nuevoCodigo"];

					if(!is_dir($absDir)){
						mkdir($absDir, 0755, true);
					}

					/*=============================================
					PROCESS IMAGE (any format supported by GD)
					=============================================*/

					$source = @imagecreatefromstring(file_get_contents($_FILES["nuevaImagen"]["tmp_name"]));

					if($source !== false){

						$random = mt_rand(100,999);

						$path = "views/img/productos/".$_POST["nuevoCodigo"]."/".$random.".jpg";

						$destination = imagecreatetruecolor($newWidth, $newHeight);

						imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);

						imagejpeg($destination, $absDir."/".$random.".jpg");

						imagedestroy($source);
						imagedestroy($destination);

					}

				}

				$table = "productos";

				$data = array("codigo" => $_POST["nuevoCodigo"],
					           "codigo_barras" => $_POST["nuevoCodigoBarras"],
					           "descripcion" => $_POST["nuevaDescripcion"],
					           "id_categoria" => $_POST["nuevaCategoria"],
					           "precio_compra" => $_POST["nuevoPrecioCompra"],
					           "precio_venta" => $_POST["nuevoPrecioVenta"],
					           "stock" => $_POST["nuevoStock"],
					           "stock_minimo" => $_POST["nuevoStockMinimo"],
					           "imagen" => $path);

				$response = ProductModel::mdlInsertProduct($table, $data);

				if($response == "ok"){

					/*=============================================
					LOG INITIAL STOCK AS INVENTORY ENTRY
					=============================================*/

					if(intval($_POST["nuevoStock"]) > 0){

						$newProduct = ProductModel::mdlShowProducts($table, "codigo", $_POST["nuevoCodigo"]);

						if($newProduct){

							InventoryModel::mdlInsertMovement("movimientos_inventario", array(
								"id_producto"   => $newProduct["id"],
								"id_usuario"    => $_SESSION["id"],
								"tipo"          => "entrada",
								"motivo"        => "Compra",
								"cantidad"      => intval($_POST["nuevoStock"]),
								"observaciones" => "Stock inicial al crear el producto",
								"id_referencia" => null
							));

						}

					}

					echo '<script>

					swal({

						type: "success",
						title: "¡El producto ha sido guardado correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "productos";

						}

					});

					</script>';

				}

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡La descripción no puede ir vacía o llevar caracteres especiales!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "productos";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	UPDATE PRODUCT
	=============================================*/

	static public function ctrUpdateProduct(){

		if(isset($_POST["editarCodigo"])){

			$table = "productos";

			/*=============================================
			CHECK FOR DUPLICATE CODE
			=============================================*/

			$existing = ProductModel::mdlShowProducts($table, "codigo", $_POST["editarCodigo"]);

			if($existing && $existing["id"] != $_POST["idProductoEditar"]){

				echo '<script>

				swal({

					type: "error",
					title: "¡Código duplicado!",
					text: "El código ingresado ya está siendo utilizado por otro producto.",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				});

				</script>';

				return;

			}

			/*=============================================
			VALIDATE IMAGE
			=============================================*/

			$path = $_POST["imagenActual"];

			if(isset($_FILES["editarImagen"]["tmp_name"]) && !empty($_FILES["editarImagen"]["tmp_name"])){

				list($width, $height) = getimagesize($_FILES["editarImagen"]["tmp_name"]);

				$newWidth = 500;
				$newHeight = 500;

				$absDir = __DIR__ . "/../views/img/productos/" . $_POST["editarCodigo"];

				if(!is_dir($absDir)){
					mkdir($absDir, 0755, true);
				}

				/*=============================================
				PROCESS IMAGE (any format supported by GD)
				=============================================*/

				$source = @imagecreatefromstring(file_get_contents($_FILES["editarImagen"]["tmp_name"]));

				if($source !== false){

					$random = mt_rand(100,999);

					$path = "views/img/productos/".$_POST["editarCodigo"]."/".$random.".jpg";

					$destination = imagecreatetruecolor($newWidth, $newHeight);

					imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);

					imagejpeg($destination, $absDir."/".$random.".jpg");

					imagedestroy($source);
					imagedestroy($destination);

				}

			}

			$data = array("id" => $_POST["idProductoEditar"],
				           "codigo" => $_POST["editarCodigo"],
				           "codigo_barras" => $_POST["editarCodigoBarras"],
				           "descripcion" => $_POST["editarDescripcion"],
				           "id_categoria" => $_POST["editarCategoria"],
				           "precio_compra" => $_POST["editarPrecioCompra"],
				           "precio_venta" => $_POST["editarPrecioVenta"],
				           "stock" => $_POST["editarStock"],
				           "stock_minimo" => $_POST["editarStockMinimo"],
				           "imagen" => $path);

			$response = ProductModel::mdlUpdateProduct($table, $data);

			if($response == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡El producto ha sido editado correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "productos";

					}

				});

				</script>';

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡Error al editar el producto!",
					text: "No se pudo actualizar el producto. Inténtelo nuevamente.",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				});

				</script>';

			}

		}

	}

	/*=============================================
	DELETE PRODUCT
	=============================================*/

	static public function ctrDeleteProduct(){

		if(isset($_GET["idProducto"])){

			$table = "productos";

			$data = array("id" => $_GET["idProducto"],
				           "estado" => 0);

			$response = ProductModel::mdlDeleteProduct($table, $data);

			if($response == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡El producto ha sido eliminado correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "productos";

					}

				});

				</script>';

			}

		}

	}

}