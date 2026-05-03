<?php

class ControladorProductos{

	/*=============================================
	MOSTRAR PRODUCTOS
	=============================================*/

	static public function ctrMostrarProductos($item, $valor){

		$tabla = "productos";

		$respuesta = ModeloProductos::mdlMostrarProductos($tabla, $item, $valor);

		return $respuesta;

	}

	/*=============================================
	CREAR PRODUCTO
	=============================================*/

	static public function ctrCrearProducto(){

		if(isset($_POST["nuevoCodigo"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevaDescripcion"])){

				/*=============================================
				VALIDAR IMAGEN
				=============================================*/

				$ruta = "";

				if(isset($_FILES["nuevaImagen"]["tmp_name"]) && !empty($_FILES["nuevaImagen"]["tmp_name"])){

					list($ancho, $alto) = getimagesize($_FILES["nuevaImagen"]["tmp_name"]);

					$nuevoAncho = 500;
					$nuevoAlto = 500;

					$directorio = "vistas/img/productos/".$_POST["nuevoCodigo"];

					mkdir($directorio, 0755);

					if($_FILES["nuevaImagen"]["type"] == "image/jpeg"){

						$aleatorio = mt_rand(100,999);
						$ruta = "vistas/img/productos/".$_POST["nuevoCodigo"]."/".$aleatorio.".jpg";
						$origen = imagecreatefromjpeg($_FILES["nuevaImagen"]["tmp_name"]);
						$destino = imagecreatetruecolor($nuevoAncho, $nuevoAlto);
						imagecopyresized($destino, $origen, 0, 0, 0, 0, $nuevoAncho, $nuevoAlto, $ancho, $alto);
						imagejpeg($destino, $ruta);

					}

					if($_FILES["nuevaImagen"]["type"] == "image/png"){

						$aleatorio = mt_rand(100,999);
						$ruta = "vistas/img/productos/".$_POST["nuevoCodigo"]."/".$aleatorio.".png";
						$origen = imagecreatefrompng($_FILES["nuevaImagen"]["tmp_name"]);
						$destino = imagecreatetruecolor($nuevoAncho, $nuevoAlto);
						imagecopyresized($destino, $origen, 0, 0, 0, 0, $nuevoAncho, $nuevoAlto, $ancho, $alto);
						imagepng($destino, $ruta);

					}

				}

				$tabla = "productos";

				$datos = array("codigo" => $_POST["nuevoCodigo"],
					           "codigo_barras" => $_POST["nuevoCodigoBarras"],
					           "descripcion" => $_POST["nuevaDescripcion"],
					           "id_categoria" => $_POST["nuevaCategoria"],
					           "precio_compra" => $_POST["nuevoPrecioCompra"],
					           "precio_venta" => $_POST["nuevoPrecioVenta"],
					           "stock" => $_POST["nuevoStock"],
					           "stock_minimo" => $_POST["nuevoStockMinimo"],
					           "imagen" => $ruta);

				$respuesta = ModeloProductos::mdlIngresarProducto($tabla, $datos);

				if($respuesta == "ok"){

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
	EDITAR PRODUCTO
	=============================================*/

	static public function ctrEditarProducto(){

		if(isset($_POST["editarCodigo"])){

			/*=============================================
			VALIDAR IMAGEN
			=============================================*/

			$ruta = $_POST["imagenActual"];

			if(isset($_FILES["editarImagen"]["tmp_name"]) && !empty($_FILES["editarImagen"]["tmp_name"])){

				list($ancho, $alto) = getimagesize($_FILES["editarImagen"]["tmp_name"]);

				$nuevoAncho = 500;
				$nuevoAlto = 500;

				$directorio = "vistas/img/productos/".$_POST["editarCodigo"];

				if(!is_dir($directorio)){
					mkdir($directorio, 0755);
				}

				if($_FILES["editarImagen"]["type"] == "image/jpeg"){

					$aleatorio = mt_rand(100,999);
					$ruta = "vistas/img/productos/".$_POST["editarCodigo"]."/".$aleatorio.".jpg";
					$origen = imagecreatefromjpeg($_FILES["editarImagen"]["tmp_name"]);
					$destino = imagecreatetruecolor($nuevoAncho, $nuevoAlto);
					imagecopyresized($destino, $origen, 0, 0, 0, 0, $nuevoAncho, $nuevoAlto, $ancho, $alto);
					imagejpeg($destino, $ruta);

				}

				if($_FILES["editarImagen"]["type"] == "image/png"){

					$aleatorio = mt_rand(100,999);
					$ruta = "vistas/img/productos/".$_POST["editarCodigo"]."/".$aleatorio.".png";
					$origen = imagecreatefrompng($_FILES["editarImagen"]["tmp_name"]);
					$destino = imagecreatetruecolor($nuevoAncho, $nuevoAlto);
					imagecopyresized($destino, $origen, 0, 0, 0, 0, $nuevoAncho, $nuevoAlto, $ancho, $alto);
					imagepng($destino, $ruta);

				}

			}

			$tabla = "productos";

			$datos = array("id" => $_POST["idProductoEditar"],
				           "codigo" => $_POST["editarCodigo"],
				           "codigo_barras" => $_POST["editarCodigoBarras"],
				           "descripcion" => $_POST["editarDescripcion"],
				           "id_categoria" => $_POST["editarCategoria"],
				           "precio_compra" => $_POST["editarPrecioCompra"],
				           "precio_venta" => $_POST["editarPrecioVenta"],
				           "stock" => $_POST["editarStock"],
				           "stock_minimo" => $_POST["editarStockMinimo"],
				           "imagen" => $ruta);

			$respuesta = ModeloProductos::mdlEditarProducto($tabla, $datos);

			if($respuesta == "ok"){

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
	BORRAR PRODUCTO
	=============================================*/

	static public function ctrBorrarProducto(){

		if(isset($_GET["idProducto"])){

			$tabla = "productos";

			$datos = array("id" => $_GET["idProducto"],
				           "estado" => 0);

			$respuesta = ModeloProductos::mdlBorrarProducto($tabla, $datos);

			if($respuesta == "ok"){

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
