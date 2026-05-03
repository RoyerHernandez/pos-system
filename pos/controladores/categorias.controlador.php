<?php

class ControladorCategorias{

	/*=============================================
	MOSTRAR CATEGORIAS
	=============================================*/

	static public function ctrMostrarCategorias($item, $valor){

		$tabla = "categorias";

		$respuesta = ModeloCategorias::mdlMostrarCategorias($tabla, $item, $valor);

		return $respuesta;

	}

	/*=============================================
	CREAR CATEGORIA
	=============================================*/

	static public function ctrCrearCategoria(){

		if(isset($_POST["nuevaCategoria"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevaCategoria"])){

				$tabla = "categorias";

				$datos = array("nombre" => $_POST["nuevaCategoria"],
					           "descripcion" => $_POST["nuevaDescripcionCat"]);

				$respuesta = ModeloCategorias::mdlIngresarCategoria($tabla, $datos);

				if($respuesta == "ok"){

					echo '<script>

					swal({

						type: "success",
						title: "¡La categoría ha sido guardada correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "categorias";

						}

					});

					</script>';

				}

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡La categoría no puede ir vacía o llevar caracteres especiales!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "categorias";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	EDITAR CATEGORIA
	=============================================*/

	static public function ctrEditarCategoria(){

		if(isset($_POST["editarCategoria"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["editarCategoria"])){

				$tabla = "categorias";

				$datos = array("id" => $_POST["idCategoriaEditar"],
					           "nombre" => $_POST["editarCategoria"],
					           "descripcion" => $_POST["editarDescripcionCat"]);

				$respuesta = ModeloCategorias::mdlEditarCategoria($tabla, $datos);

				if($respuesta == "ok"){

					echo '<script>

					swal({

						type: "success",
						title: "¡La categoría ha sido editada correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "categorias";

						}

					});

					</script>';

				}

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡La categoría no puede ir vacía o llevar caracteres especiales!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "categorias";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	BORRAR CATEGORIA
	=============================================*/

	static public function ctrBorrarCategoria(){

		if(isset($_GET["idCategoria"])){

			$tabla = "categorias";

			$datos = array("id" => $_GET["idCategoria"],
				           "estado" => 0);

			$respuesta = ModeloCategorias::mdlBorrarCategoria($tabla, $datos);

			if($respuesta == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡La categoría ha sido eliminada correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "categorias";

					}

				});

				</script>';

			}

		}

	}

}
