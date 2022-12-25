<script>
export default {
	data: function() {
		return {
			errormsg: null,
			searchName: "",
			baseName: "",
			result: ""
		}
	},
	methods: {
		async userSearch() {
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", { params: {
                        "user-name": this.searchName
                    }});
				this.baseName = response.data[response.data.length - 1];
				this.result = response.data;
				this.emptyDropdown();
				this.populateDropdown();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		resetSearch() {
			this.baseName = "";
			this.result = "";
			this.emptyDropdown();
		},

		populateDropdown() {
			var dropdown = document.getElementById("candidateUsersDropdown");

			this.result.forEach(element => {
				var but = document.createElement("button");
				but.addEventListener('click', () => this.$router.push({name: 'user', params: {username: element}}));
				var text = document.createTextNode(element);
				but.classList.add("navbar-search-result")
				but.appendChild(text);
				dropdown.appendChild(but);
			});

			if (this.result.length >= 64) {
				var but = document.createElement("button");
				but.addEventListener('click', () => this.furtherRequest());
				but.setAttribute("id", "furtherRequest_button")
				var text = document.createTextNode("Click here to see more results!");
				but.classList.add("navbar-search-result-end")
				but.appendChild(text);
				dropdown.appendChild(but);
			}
		},

		async furtherRequest() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", { params: {
                        "user-name": this.searchName, "name-base": this.baseName
                    }});
				this.baseName = response.data[response.data.length - 1];
				this.result = response.data;
				document.getElementById("furtherRequest_button").remove();
				this.populateDropdown();
			} catch (e) {
				this.errormsg = e.toString();
				this.emptyDropdown();
			}
			this.loading = false;
		},

		emptyDropdown() {
			var dropdown = document.getElementById("candidateUsersDropdown");
			while (dropdown.firstChild) {
				dropdown.removeChild(dropdown.lastChild);
			}
		}
	}
}
</script>

<template>
    <div class="navbar-search">
        <input class="navbar-search-input text-white col-lg-1 px-4 fs-6" v-on:keyup="resetSearch" v-on:keyup.enter="userSearch" v-model="searchName" type="text" placeholder="Search user">
        <div class="dropdown-menu-dark disable-scrollbars navbar-search-results no-bullets" id="candidateUsersDropdown" aria-labelledby="searchDropdown">
        </div>
    </div>
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>