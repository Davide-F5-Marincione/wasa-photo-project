<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar navbar-dark sticky-top flex-md-nowrap p-0">
		<a class="navbar-brand col-md-3 col-lg-1 me-0 px-4 fs-5" href="#/">WASAPhoto</a>
		<input class="navbar-search text-white col-lg-1 px-4 fs-6" placeholder="Search user">
	</header>

	<div class="container-fluid">
		<div class="row">
			<main class="">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
