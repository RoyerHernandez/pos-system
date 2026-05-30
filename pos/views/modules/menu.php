<?php $role = isset($_SESSION["role"]) ? $_SESSION["role"] : "Vendedor"; ?>

<aside class="main-sidebar">

	 <section class="sidebar">

		<ul class="sidebar-menu">

			<!-- INICIO - All roles -->
			<li class="active">

				<a href="inicio">

					<i class="fa fa-home"></i>
					<span>Inicio</span>

				</a>

			</li>

			<?php if($role == "Administrador"): ?>
			<!-- USUARIOS - Admin only -->
			<li>

				<a href="usuarios">

					<i class="fa fa-user"></i>
					<span>Usuarios</span>

				</a>

			</li>
			<?php endif; ?>

			<?php if($role == "Administrador"): ?>
			<!-- CATEGORIAS - Admin only -->
			<li>

				<a href="categorias">

					<i class="fa fa-th"></i>
					<span>Categorias</span>

				</a>

			</li>
			<?php endif; ?>

			<?php if($role == "Administrador" || $role == "Especial"): ?>
			<!-- PRODUCTOS - Admin + Especial -->
			<li>

				<a href="productos">

					<i class="fa fa-product-hunt"></i>
					<span>Productos</span>

				</a>

			</li>
			<?php endif; ?>

			<?php if($role == "Administrador" || $role == "Especial"): ?>
			<!-- CLIENTES - Admin + Especial -->
			<li>

				<a href="clientes">

					<i class="fa fa-users"></i>
					<span>Clientes</span>

				</a>

			</li>
			<?php endif; ?>

			<!-- CAJA - All roles -->
			<li>

				<a href="caja">

					<i class="fa fa-desktop"></i>
					<span>Caja</span>

				</a>

			</li>

			<!-- VENTAS - Submenu -->
			<li class="treeview">

				<a href="#">

					<i class="fa fa-list-ul"></i>

					<span>Ventas</span>

					<span class="pull-right-container">

						<i class="fa fa-angle-left pull-right"></i>

					</span>

				</a>

				<ul class="treeview-menu">

					<?php if($role == "Administrador" || $role == "Especial"): ?>
					<!-- ADMINISTRAR VENTAS - Admin + Especial -->
					<li>

						<a href="ventas">

							<i class="fa fa-circle-o"></i>
							<span>Administrar ventas</span>

						</a>

					</li>
					<?php endif; ?>

					<!-- CREAR VENTA - All roles -->
					<li>

						<a href="crear-venta">

							<i class="fa fa-circle-o"></i>
							<span>Crear venta</span>

						</a>

					</li>

					<?php if($role == "Administrador" || $role == "Especial"): ?>
					<!-- REPORTES - Admin + Especial -->
					<li>

						<a href="reportes">

							<i class="fa fa-circle-o"></i>
							<span>Reporte de ventas</span>

						</a>

					</li>
					<?php endif; ?>

				</ul>

			</li>

		</ul>

	 </section>

</aside>
