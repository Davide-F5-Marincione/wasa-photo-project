<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
            input: "",
            cohort: "",
            loginhandle: "",
		}
	},
	methods: {
		async login() {
            this.loginhandle = this.input + "#" + this.cohort

            this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/users", null, { params: {
                        "user-handle": this.loginhandle
                    }});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;

		},
	}
}
</script>

<template>
	<div>
        <div class="justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3">
            <h5 class="h5">Login:</h5><br>
            Handle: <input v-on:keyup.enter="login" v-model="input" type="text" placeholder="handle's name" />
            <h7 class="h7">#</h7><input v-on:keyup.enter="login"  v-model="cohort" type="text" placeholder="0000"/>
            <p> The response is: <b>{{some_data}}</b> </p>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>