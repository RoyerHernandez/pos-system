<div class="content-wrapper">

  <section class="content-header">
    
    <h1>
      
      Administrar usuarios
    
    </h1>

    <ol class="breadcrumb">
      
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      
      <li class="active">Administrar usuarios</li>
    
    </ol>

  </section>

  <!-- Main content -->
  <section class="content">

    <!-- Default box -->
    <div class="box">

      <div class="box-header with-border">

        <button class="btn btn-primary" data-toggle="modal" data-target="#modalAgregarUsuario">
          Agregar Usuario
        </button> 
       
      </div>

      <div class="box-body">
        
        <table class="table table-bordered table-striped dt-responsive tablas">
          
          <thead>
            
          <th style="width: 10px"></th>
            <th>Nombre</th>
            <th>Usuario</th>
            <th>Foto</th>
            <th>Perfil</th>
            <th>Estado</th>
            <th>Último login</th>
            <th>Acciones</th>

          </thead>

          <tbody>
            
            <tr>
              
            <td>1</td>
            <td>Usuario Administrador</td>
            <td>admin</td>
            <td><img src="vistas/img/usuarios/default/anonymous.png" class="img-thumbnail" width="40px"></td>
            <td>Administrador</td>
            <td><button class="btn btn-success btn-xs">Activado</button></td>
            <td>2024-12-11 12::32</td>
            <td>
              
              <div class="btn-group">
                
                <button class="btn btn-warning"><i class="fa fa-pencil"></i></button>

                <button class="btn btn-danger"><i class="fa fa-times"></i></button>

              </div>

            </td>
 

            </tr>
            <tr>
              
            <td>1</td>
            <td>Usuario Administrador</td>
            <td>admin</td>
            <td><img src="vistas/img/usuarios/default/anonymous.png" class="img-thumbnail" width="40px"></td>
            <td>Administrador</td>
            <td><button class="btn btn-success btn-xs">Activado</button></td>
            <td>2024-12-11 12::32</td>
            <td>
              
              <div class="btn-group">
                
                <button class="btn btn-warning"><i class="fa fa-pencil"></i></button>

                <button class="btn btn-danger"><i class="fa fa-times"></i></button>

              </div>

            </td>
 

            </tr>
            <tr>
              
            <td>1</td>
            <td>Usuario Administrador</td>
            <td>admin</td>
            <td><img src="vistas/img/usuarios/default/anonymous.png" class="img-thumbnail" width="40px"></td>
            <td>Administrador</td>
            <td><button class="btn btn-danger btn-xs">Inactivo</button></td>
            <td>2024-12-11 12::32</td>
            <td>
              
              <div class="btn-group">
                
                <button class="btn btn-warning"><i class="fa fa-pencil"></i></button>

                <button class="btn btn-danger"><i class="fa fa-times"></i></button>

              </div>

            </td>
 

            </tr>

          </tbody>

        </table>

      </div>
      
    </div>

  </section>
  
</div>

<!--Modal Agregar Usuario-->
<div class="modal fade" id="modalAgregarUsuario">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post" enctype="multipart/form-data">  

          <!-- Modal Header -->
          
          <div class="modal-header" style="background:#3c8dbc; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times</button>
            <h4 class="modal-title">Agregar Usuario</h4>
         
          </div>
          
          <div class="modal-body">
            

              <div class="box-body">

                <!--ENTRADA PARA EL NOMBRE-->

                  <div class="form-group">
                    
                    <div class="input-group">

                      <span class="input-group-addon"><i class="fa  fa-user"></i></span>

                      <input type="text" class="form-control input-lg" name="nuevoNombre" placeholder="Ingresar nombre" required>
                      
                    </div>

                  </div>

                  <!--ENTRADA PARA EL USUARIO-->

                  <div class="form-group">
                    
                    <div class="input-group">

                      <span class="input-group-addon"><i class="fa  fa-key"></i></span>

                      <input type="text" class="form-control input-lg" name="nuevoUsuario" placeholder="Ingresar usuario" required>
                      
                    </div>

                  </div>

                  <!--ENTRADA PARA EL CONTRASEÑA-->

                  <div class="form-group">
                    
                    <div class="input-group">

                      <span class="input-group-addon"><i class="fa  fa-lock"></i></span>

                      <input type="password" class="form-control input-lg" name="nuevoPassword" placeholder="Ingresar contraseña" required>
                      
                    </div>

                  </div>

                  <!--ENTRADA PARA SELECCIONA SU PERFIL-->

                  <div class="form-group">
                    
                    <div class="input-group">

                      <span class="input-group-addon"><i class="fa  fa-user"></i></span>

                     <select class="form-control input-lg" name="nuevoPerfil">
                       
                       <option value="">Seleccionar perfil</option>

                       <option value="Administrado">Administrado</option>

                       <option value="Especial">Especial</option>

                       <option value="Vendedor">Vendedoer</option>                

                     </select>
                      
                    </div>

                  </div>

                  <!--ENTRADA PARA SUBIR FOTO-->

                  <div class="form-group">
                    
                    <div class="panel">SUBIR FOTO</div>

                    <p class="help-block">peso maximo de la foto 200 MB</p>

                    <img src=" vistas/img/usuarios/default/anonymous.png" class=" img-thumbnail" width="100px">

                  </div>

              </div>

          </div>

          <div class="modal-footer">

            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>

            <button type="submint" class="btn btn-primary">Guardar usuario</button>

          </div>

      </form>

    </div>
  </div>
</div>