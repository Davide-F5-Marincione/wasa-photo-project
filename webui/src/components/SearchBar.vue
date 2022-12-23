<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			searchName: "",
			baseName: "",
			result: "",
		}
	},
	methods: {
		async userSearch() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", { params: {
                        "user-name": this.searchName
                    }});
				this.baseName = response.data[-1]
				this.result = response.data
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		resetSearch() {
			this.baseName = ""
			this.result = ""
		},
	}
}
</script>

<template>
    <div class="navbar-search">
        <input class="navbar-search-input text-white col-lg-1 px-4 fs-6" v-on:keyup="resetSearch" v-on:keyup.enter="userSearch" v-model="searchName" type="text" placeholder="Search user">
        <ul class="dropdown-menu-dark navbar-search-results" aria-labelledby="searchDropdown">
            <li><a class="dropdown-item" href="#">User 1</a></li>
            <li><a class="dropdown-item" href="#">User 2</a></li>
            <li><a class="dropdown-item" href="#">User 3</a></li>
        </ul>
    </div>
</template>

<style>
</style>