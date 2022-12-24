<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			username: "",
			token: 0
		}
	},

	mounted() {
		this.username = localStorage.username;
		this.token = localStorage.token;
	}
}
</script>

<template>
	<header class="navbar navbar-dark sticky-top flex-md-nowrap p-0">
		<a class="navbar-brand col-md-3 col-lg-1 me-0 px-4 fs-5" href="#/">WASAPhoto</a>
		<SearchBar></SearchBar>
		<div class="text-white fs-6">Logged as: <b class="col-lg-1 px-10">{{username}}</b></div>
	</header>
</template>

<style>
</style>
