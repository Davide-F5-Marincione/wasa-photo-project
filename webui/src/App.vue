<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<main class="">
				<!-- The matter that when pushing to the same route but with different parameters
					is a behaviour not supported by default (and therefore one has to add this
					sort of hacky-stuff) is bullshit. -->
				<RouterView :key="$route.fullPath"/>
			</main>
		</div>
	</div>
</template>

<style>
</style>
