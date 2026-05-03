 <header class="main-header">

	<!--=====================================
	LOGOTIPO
	======================================-->
	<a href="inicio" class="logo">

		<!-- logo mini -->
		<span class="logo-mini">

			<img src="vistas/img/plantilla/icono-blanco.png" class="img-responsive" style="padding:10px">

		</span>

		<!-- logo normal -->

		<span class="logo-lg">

			<img src="vistas/img/plantilla/logo-blanco-lineal.png" class="img-responsive" style="padding:10px 0px">

		</span>

	</a>

	<!--=====================================
	BARRA DE NAVEGACIÓN
	======================================-->
	<nav class="navbar navbar-static-top" role="navigation">

		<!-- Botón de navegación -->

	 	<a href="#" class="sidebar-toggle" data-toggle="push-menu" role="button">

        	<span class="sr-only">Toggle navigation</span>

      	</a>

		<!-- perfil de usuario -->

		<div class="navbar-custom-menu">

			<ul class="nav navbar-nav">

				<li class="dropdown user user-menu">

					<a href="#" class="dropdown-toggle" data-toggle="dropdown">

						<?php if(!empty($_SESSION["foto"])): ?>
							<img src="<?php echo $_SESSION["foto"]; ?>" class="user-image">
						<?php else: ?>
							<img src="vistas/img/usuarios/default/anonymous.png" class="user-image">
						<?php endif; ?>

						<span class="hidden-xs"><?php echo isset($_SESSION["nombre"]) ? $_SESSION["nombre"] : "Usuario"; ?></span>

					</a>

					<!-- Dropdown-toggle -->

					<ul class="dropdown-menu">

						<!-- User image -->
						<li class="user-header" style="background:#3c8dbc;">

							<?php if(!empty($_SESSION["foto"])): ?>
								<img src="<?php echo $_SESSION["foto"]; ?>" class="img-circle" alt="User Image">
							<?php else: ?>
								<img src="vistas/img/usuarios/default/anonymous.png" class="img-circle" alt="User Image">
							<?php endif; ?>

							<p>
								<?php echo isset($_SESSION["nombre"]) ? $_SESSION["nombre"] : "Usuario"; ?>
								<small><?php echo isset($_SESSION["perfil"]) ? $_SESSION["perfil"] : ""; ?></small>
							</p>

						</li>

						<li class="user-body">

							<div class="pull-right">

								<a href="salir" class="btn btn-default btn-flat">Salir</a>

							</div>

						</li>

					</ul>

				</li>

			</ul>

		</div>

	</nav>

 </header>