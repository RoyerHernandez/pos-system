<?php

$categorias = ControladorCategorias::ctrMostrarCategorias(null, null);

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

            $contador = 1;

            foreach($categorias as $key => $valor){

              echo '<tr>

                <td>'.$contador.'</td>
                <td>'.$valor["nombre"].'</td>
                <td>'.$valor["descripcion"].'</td>
                <td>';

              if($valor["estado"] == 1){
                echo '<button class="btn btn-success btn-xs">Activada</button>';
              }else{
                echo '<button class="btn btn-danger btn-xs">Inactiva</button>';
              }

              echo '</td>
                <td>

                  <div class="btn-group">

                    <button class="btn btn-warning btnEditarCategoria" data-toggle="modal" data-target="#modalEditarCategoria" idCategoria="'.$valor["id"].'" catNombre="'.$valor["nombre"].'" catDescripcion="'.$valor["descripcion"].'"><i class="fa fa-pencil"></i></button>

                    <a class="btn btn-danger btnEliminarCategoria" idCategoria="'.$valor["id"].'"><i class="fa fa-times"></i></a>

                  </div>

                </td>

              </tr>';

              $contador++;

            }

            ?>

          </tbody>

        </table>

      </div>

    </div>

  </section>

</div>

<!--=====================================
MODAL AGREGAR CATEGORIA
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

                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevaCategoria" placeholder="Nombre de la categoría" required>
                  </div>
                </div>

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
MODAL EDITAR CATEGORIA
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

                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <input type="text" class="form-control input-lg" name="editarCategoria" id="editarCategoria" placeholder="Nombre de la categoría" required>
                  </div>
                </div>

                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-comment"></i></span>
                    <input type="text" class="form-control input-lg" name="editarDescripcionCat" id="editarDescripcionCat" placeholder="Descripción (opcional)">
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

$crearCategoria = new ControladorCategorias();
$crearCategoria -> ctrCrearCategoria();

$editarCategoria = new ControladorCategorias();
$editarCategoria -> ctrEditarCategoria();

$borrarCategoria = new ControladorCategorias();
$borrarCategoria -> ctrBorrarCategoria();

?>

<script>

/*=============================================
EDITAR CATEGORIA - Cargar datos en modal
=============================================*/

$(".tablas").on("click", ".btnEditarCategoria", function(){

  var idCategoria = $(this).attr("idCategoria");
  var catNombre = $(this).attr("catNombre");
  var catDescripcion = $(this).attr("catDescripcion");

  $("#idCategoriaEditar").val(idCategoria);
  $("#editarCategoria").val(catNombre);
  $("#editarDescripcionCat").val(catDescripcion);

});

/*=============================================
ELIMINAR CATEGORIA
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
