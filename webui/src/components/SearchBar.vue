<script>
export default {
	data: function() {
		return {
			errormsg: null,
			searchName: "",
			baseName: "",
			results: []
		}
	},
	methods: {
		async userSearch() {
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", { params: {
                        "user-name": this.searchName, "name-base": this.baseName
                    }});
				this.results.push(...response.data);
				if (response.data.length > 0) {
					this.baseName = response.data[response.data.length - 1];
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		resetSearch() {
			this.baseName = "";
			this.results = [];
		}
	}
}
</script>

<template>
    <div class="navbar-search">
        <input class="navbar-search-input" v-on:keyup="resetSearch" v-on:keyup.enter="userSearch" v-model="searchName" type="text" placeholder="Search user">
        <div class="dropdown-menu-dark disable-scrollbars" id="candidateUsersDropdown">
			<div v-for="element in results"  class="navbar-search-result">
				<router-link class="navbar-search-result-text" v-bind:to="'/users/' + element"> {{ element }}</router-link>
			</div>
			<button v-if="results.length > 0" class="navbar-search-result-end" :onclick="() => this.furtherRequest()">Click here to see more results!</button>
        </div>
    </div>
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>