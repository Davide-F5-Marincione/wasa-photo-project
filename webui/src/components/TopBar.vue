<script>
import SearchBar from '../components/SearchBar.vue';

export default {
	data: function() {
		return {
			username: "",
			token: 0
		}
	},
	mounted() {
		this.username = localStorage.username;
		this.token = localStorage.token;
	},
    components: {SearchBar}
}
</script>

<template>
	<header class="navbar navbar-dark my-topbar">
		<router-link class="navbar-brand" to="/">WASAPhoto</router-link>
        <router-link class="navbar-goto-stream" to="/stream">Goto stream</router-link>
		<router-link class="navbar-goto-upload" to="/upload">Upload photo</router-link>
		<SearchBar></SearchBar>
		<div class="logged-text">Logged as: <router-link class="logged-text" :to="'/users/' + username"><b>{{username}}</b></router-link></div>
	</header>
</template>

<style>
</style>
