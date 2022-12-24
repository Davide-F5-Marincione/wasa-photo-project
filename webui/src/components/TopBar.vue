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
	<header class="navbar navbar-dark sticky-top flex-md-nowrap p-0">
		<a class="navbar-brand col-md-3 col-lg-1 me-0 px-4 fs-5" href="#/">WASAPhoto</a>
        <button class="navbar-goto-stream col-md-3 col-lg-1 me-0 px-4 fs-7" @click="toStream()">
            <a>Goto stream</a>
        </button>
		<SearchBar></SearchBar>
		<div class="text-white fs-6">Logged as: <b class="col-lg-1 px-10">{{username}}</b></div>
	</header>
</template>

<style>
</style>
