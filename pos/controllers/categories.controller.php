<?php

class CategoryController{

	/*=============================================
	SHOW CATEGORIES
	=============================================*/

	static public function ctrShowCategories($item, $value){

		$table = "categorias";

		$response = CategoryModel::mdlShowCategories($table, $item, $value);

		return $response;

	}

	/*=============================================
	CREATE CATEGORY
	=============================================*/

	static public function ctrCreateCategory(){

		if(isset($_POST["nuevaCategoria"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevaCategoria"])){

				$table = "categorias";

				$data = array("nombre" => $_POST["nuevaCategoria"],
					           "descripcion" => $_POST["nuevaDescripcionCat"]);

				$response = CategoryModel::mdlInsertCategory($table, $data);

				if($response == "ok"){

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
	UPDATE CATEGORY
	=============================================*/

	static public function ctrUpdateCategory(){

		if(isset($_POST["editarCategoria"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["editarCategoria"])){

				$table = "categorias";

				$data = array("id" => $_POST["idCategoriaEditar"],
					           "nombre" => $_POST["editarCategoria"],
					           "descripcion" => $_POST["editarDescripcionCat"]);

				$response = CategoryModel::mdlUpdateCategory($table, $data);

				if($response == "ok"){

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
	DELETE CATEGORY
	=============================================*/

	static public function ctrDeleteCategory(){

		if(isset($_GET["idCategoria"])){

			$table = "categorias";

			$data = array("id" => $_GET["idCategoria"],
				           "estado" => 0);

			$response = CategoryModel::mdlDeleteCategory($table, $data);

			if($response == "ok"){

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