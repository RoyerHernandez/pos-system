 <header class="main-header">

	<!--=====================================
	LOGO
	======================================-->
	<a href="inicio" class="logo">

		<!-- logo mini -->
		<span class="logo-mini">

			<img src="views/img/plantilla/icono-blanco.svg" class="img-responsive" style="padding:10px">

		</span>

		<!-- logo normal -->

		<span class="logo-lg">

			<img src="views/img/plantilla/logo-callejon-blanco-lineal.svg" class="img-responsive" style="padding:10px 0px">

		</span>

	</a>

	<!--=====================================
	NAVIGATION BAR
	======================================-->
	<nav class="navbar navbar-static-top" role="navigation">

		<!-- Navigation button -->

	 	<a href="#" class="sidebar-toggle" data-toggle="push-menu" role="button">

        	<span class="sr-only">Toggle navigation</span>

      	</a>

		<!-- User profile -->

		<div class="navbar-custom-menu">

			<ul class="nav navbar-nav">

				<li class="dropdown user user-menu">

					<a href="#" class="dropdown-toggle" data-toggle="dropdown">

						<?php if(!empty($_SESSION["photo"])): ?>
							<img src="<?php echo $_SESSION["photo"]; ?>" class="user-image">
						<?php else: ?>
							<img src="views/img/usuarios/default/avatar-default.svg" class="user-image">
						<?php endif; ?>

						<span class="hidden-xs"><?php echo isset($_SESSION["name"]) ? $_SESSION["name"] : "Usuario"; ?></span>

					</a>

					<!-- Dropdown-toggle -->

					<ul class="dropdown-menu">

						<!-- User image -->
						<li class="user-header" style="background:#2C2C2C;">

							<?php if(!empty($_SESSION["photo"])): ?>
								<img src="<?php echo $_SESSION["photo"]; ?>" class="img-circle" alt="User Image">
							<?php else: ?>
								<img src="views/img/usuarios/default/avatar-default.svg" class="img-circle" alt="User Image">
							<?php endif; ?>

							<p>
								<?php echo isset($_SESSION["name"]) ? $_SESSION["name"] : "Usuario"; ?>
								<small><?php echo isset($_SESSION["role"]) ? $_SESSION["role"] : ""; ?></small>
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