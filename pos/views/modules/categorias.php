<?php

$categories = CategoryController::ctrShowCategories(null, null);

?>

<div class="content-wrapper">

  <section class="content-header">

    <h1>
      Administrar categorías
    </h1>

    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Categorías</li>
    </ol>

  </section>

  <section class="content">

    <div class="box">

      <div class="box-header with-border">

        <button class="btn btn-primary" data-toggle="modal" data-target="#modalAgregarCategoria">
          <i class="fa fa-plus"></i> Agregar Categoría
        </button>

      </div>

      <div class="box-body">

        <table class="table table-bordered table-striped dt-responsive tablas">

          <thead>
            <tr>
              <th style="width:10px;">#</th>
              <th>Nombre</th>
              <th>Descripción</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>

          <tbody>

            <?php

            $counter = 1;

            foreach($categories as $key => $value){

              echo '<tr>

                <td>'.$counter.'</td>
                <td>'.$value["nombre"].'</td>
                <td>'.$value["descripcion"].'</td>
                <td>';

              if($value["estado"] == 1){
                echo '<button class="btn btn-success btn-xs">Activado</button>';
              }else{
                echo '<button class="btn btn-danger btn-xs">Inactivo</button>';
              }

              echo '</td>
                <td>

                  <div class="btn-group">

                    <button class="btn btn-warning btnEditarCategoria" data-toggle="modal" data-target="#modalEditarCategoria" idCategoria="'.$value["id"].'" nombre="'.$value["nombre"].'" descripcion="'.$value["descripcion"].'"><i class="fa fa-pencil"></i></button>

                    <a class="btn btn-danger btnEliminarCategoria" idCategoria="'.$value["id"].'"><i class="fa fa-times"></i></a>

                  </div>

                </td>

              </tr>';

              $counter++;

            }

            ?>

          </tbody>

        </table>

      </div>

    </div>

  </section>

</div>

<!--=====================================
ADD CATEGORY MODAL
======================================-->

<div class="modal fade" id="modalAgregarCategoria">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post">

          <div class="modal-header" style="background:#3c8dbc; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Agregar Categoría</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <!-- NAME INPUT -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevaCategoria" placeholder="Nombre de la categoría" required>
                  </div>
                </div>

                <!-- DESCRIPTION INPUT -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-comment"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevaDescripcionCat" placeholder="Descripción (opcional)">
                  </div>
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-primary">Guardar categoría</button>
          </div>

      </form>

    </div>
  </div>
</div>

<!--=====================================
EDIT CATEGORY MODAL
======================================-->

<div class="modal fade" id="modalEditarCategoria">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post">

          <div class="modal-header" style="background:#f39c12; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Editar Categoría</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <input type="hidden" name="idCategoriaEditar" id="idCategoriaEditar">

                <!-- NAME INPUT -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <input type="text" class="form-control input-lg" name="editarCategoria" id="editarCategoria" placeholder="Nombre" required>
                  </div>
                </div>

                <!-- DESCRIPTION INPUT -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-comment"></i></span>
                    <input type="text" class="form-control input-lg" name="editarDescripcionCat" id="editarDescripcionCat" placeholder="Descripción">
                  </div>
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-warning">Editar categoría</button>
          </div>

      </form>

    </div>
  </div>
</div>

<?php

$createCategory = new CategoryController();
$createCategory -> ctrCreateCategory();

$updateCategory = new CategoryController();
$updateCategory -> ctrUpdateCategory();

$deleteCategory = new CategoryController();
$deleteCategory -> ctrDeleteCategory();

?>

<script>

/*=============================================
EDIT CATEGORY - Load data into modal
=============================================*/

$(".tablas").on("click", ".btnEditarCategoria", function(){

  var idCategoria = $(this).attr("idCategoria");
  var nombre = $(this).attr("nombre");
  var descripcion = $(this).attr("descripcion");

  $("#idCategoriaEditar").val(idCategoria);
  $("#editarCategoria").val(nombre);
  $("#editarDescripcionCat").val(descripcion);

});

/*=============================================
DELETE CATEGORY
=============================================*/

$(".tablas").on("click", ".btnEliminarCategoria", function(){

  var idCategoria = $(this).attr("idCategoria");

  swal({

    title: "¿Está seguro de eliminar la categoría?",
    text: "¡Si no lo está, puede cancelar la acción!",
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#DD6B55",
    cancelButtonText: "Cancelar",
    confirmButtonText: "¡Sí, eliminar!",
    closeOnConfirm: false

  }).then(function(result){

    if(result.value){

      window.location = "categorias&idCategoria="+idCategoria;

    }

  });

});

</script>
