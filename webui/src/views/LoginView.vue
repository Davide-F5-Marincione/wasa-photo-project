<script>
import { inject } from 'vue';

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
            input: ""
		}
	},
	methods: {
		async login() {
            this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/users", null, { params: {
                        "user-name": this.input
                    }});
				this.$username.value = response.data["resp-username"]
				this.$token.value = response.data["resp-authtoken"]
				this.$router.push('/stream');
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	}
}
</script>

<template>
	<div class="navbar navbar-dark">
		<div class="navbar-logo px-4 fs-1">WASAPhoto</div>
	</div>
	<div class="center">
    	<p><input class="login-bar text-white text-center fs-4" v-on:keyup.enter="login" v-model="input" type="text" placeholder="Enter your username"/></p>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>