<script>
export default {
	data: function() {
		return {
			errormsg: null,
            input: ""
		}
	},
	methods: {
		async login() {
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/users", this.input, {headers: {"Content-Type": "application/json"}});
				localStorage.username = response.data["resp-username"];
				localStorage.token = response.data["resp-authtoken"].toString();
				this.$router.push('/stream');
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	},
	created() {
		localStorage.clear() //If we are accessing here that means we can clear the localstorage ()
	}
}
</script>

<template>
	<div class="navbar navbar-dark">
		<div class="logo px-4 fs-1">WASAPhoto</div>
	</div>
	<div class="center">
    	<p><input class="login-bar text-white text-center fs-4" v-on:keyup.enter="login" v-model="input" type="text" placeholder="Enter your username" maxlength="32"/></p>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>