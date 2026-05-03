<?php

class ProductController{

	/*=============================================
	SHOW PRODUCTS
	=============================================*/

	static public function ctrShowProducts($item, $valor){

		$table = "productos";

		$response = ProductModel::mdlShowProducts($table, $item, $valor);

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

					$directory = "vistas/img/productos/".$_POST["nuevoCodigo"];

					mkdir($directory, 0755);

					if($_FILES["nuevaImagen"]["type"] == "image/jpeg"){

						$random = mt_rand(100,999);
						$path = "vistas/img/productos/".$_POST["nuevoCodigo"]."/".$random.".jpg";
						$source = imagecreatefromjpeg($_FILES["nuevaImagen"]["tmp_name"]);
						$destination = imagecreatetruecolor($newWidth, $newHeight);
						imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);
						imagejpeg($destination, $path);

					}

					if($_FILES["nuevaImagen"]["type"] == "image/png"){

						$random = mt_rand(100,999);
						$path = "vistas/img/productos/".$_POST["nuevoCodigo"]."/".$random.".png";
						$source = imagecreatefrompng($_FILES["nuevaImagen"]["tmp_name"]);
						$destination = imagecreatetruecolor($newWidth, $newHeight);
						imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);
						imagepng($destination, $path);

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

			/*=============================================
			VALIDATE IMAGE
			=============================================*/

			$path = $_POST["imagenActual"];

			if(isset($_FILES["editarImagen"]["tmp_name"]) && !empty($_FILES["editarImagen"]["tmp_name"])){

				list($width, $height) = getimagesize($_FILES["editarImagen"]["tmp_name"]);

				$newWidth = 500;
				$newHeight = 500;

				$directory = "vistas/img/productos/".$_POST["editarCodigo"];

				if(!is_dir($directory)){
					mkdir($directory, 0755);
				}

				if($_FILES["editarImagen"]["type"] == "image/jpeg"){

					$random = mt_rand(100,999);
					$path = "vistas/img/productos/".$_POST["editarCodigo"]."/".$random.".jpg";
					$source = imagecreatefromjpeg($_FILES["editarImagen"]["tmp_name"]);
					$destination = imagecreatetruecolor($newWidth, $newHeight);
					imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);
					imagejpeg($destination, $path);

				}

				if($_FILES["editarImagen"]["type"] == "image/png"){

					$random = mt_rand(100,999);
					$path = "vistas/img/productos/".$_POST["editarCodigo"]."/".$random.".png";
					$source = imagecreatefrompng($_FILES["editarImagen"]["tmp_name"]);
					$destination = imagecreatetruecolor($newWidth, $newHeight);
					imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);
					imagepng($destination, $path);

				}

			}

			$table = "productos";

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