<script>
import SearchBar from '../components/SearchBar.vue';

export default {
	data: function() {
		return {
			username: "",
			token: 0
		}
	},

    methods: {
        toStream() {
            this.$router.push({ name: "stream"})
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
	<header class="navbar navbar-dark sticky-top my-bg-dark">
		<a class="navbar-brand" href="#/">WASAPhoto</a>
        <button class="navbar-goto-stream" @click="toStream()">
            <a>Goto stream</a>
        </button>
		<SearchBar></SearchBar>
		<div class="logged-text">Logged as: <b>{{username}}</b></div>
	</header>
</template>

<style>
</style>
