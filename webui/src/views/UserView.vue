
<script>
export default {
    data: function () {
        return {
            errormsg: null,
            photosResults: [],
			followersResults: [],
			followingResults: [],
			photosLimit: "",
			followersBase: "",
			followingBase: "",
			following: false,
			banning: false,
			username: localStorage.username,
			othername: this.$route.params.username,
			input: ""
        };
    },
    methods: {
		async checkCore() {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/users/" + this.username + "/follows/" + this.othername);
				this.following = response.data;
				var response = await this.$axios.get("/users/" + this.username + "/bans/" + this.othername);
				this.banning = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
        async refreshData(photos_new=false, followers_new=false, following_new=false) {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/users/" + this.othername, {params: {"photos-limit":this.photosLimit, "followers-base": this.followersBase, "following-base": this.followingBase}});
				
				for (let i=0; i< response.data["followers-running-batch"].length; i++) {
					if (response.data["followers-running-batch"][i].name == this.username) {
						response.data["followers-running-batch"].splice(i, 1)
						break
					}
				}

				if (photos_new) {
					this.photosResults.push(...response.data["photos-running-batch"])
					if (response.data["photos-running-batch"].length > 0) {
						this.photosLimit = response.data["photos-running-batch"][response.data["photos-running-batch"].length - 1];
					}
				}

				if (followers_new) {
					this.followersResults.push(...response.data["followers-running-batch"])
					if (response.data["followers-running-batch"].length > 0) {
						this.followersBase = response.data["followers-running-batch"][response.data["followers-running-batch"].length - 1].name;
					}
				}

				if (following_new) {
					this.followingResults.push(...response.data["following-running-batch"])
					if (response.data["following-running-batch"].length > 0) {
						this.followingBase = response.data["following-running-batch"][response.data["following-running-batch"].length - 1].name;
					}
				}

			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async banThis() {
			try {
				await this.$axios.put("/users/" + this.username + "/bans/" + this.othername);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.checkCore();
		},
		async unbanThis() {
			try {
				await this.$axios.delete("/users/" + this.username + "/bans/" + this.othername);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.checkCore();
		},
		async followThis() {
			try {
				await this.$axios.put("/users/" + this.username + "/follows/" + this.othername);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.checkCore();
		},
		async unfollowThis() {
			try {
				await this.$axios.delete("/users/" + this.username + "/follows/" + this.othername);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.checkCore();
		},
        async delPost(id) {
			this.errormsg = null;
			try {
				await this.$axios.delete("/photos/" + id.toString());
                this.results = []
                this.limit = ""
                this.refreshData()
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async newUsername() {
			this.errormsg = null;
			try {
				await this.$axios.put("/users/" + this.username, this.input, {headers:{"Content-Type":"application/json"}});
				localStorage.username = this.input
				this.input=""
                this.$router.push("/users/" + localStorage.username)
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
    },
    created() {
		this.checkCore();
        this.refreshData(true, true, true);
    }
}
</script>

<template>
	<TopBar></TopBar>
	<div>
		<div class="deck">
			<div>Checking <span v-if="othername != username"><b>{{ othername }}</b>'s</span><span v-if="othername == username"><b>your</b></span> profile.</div>
			<div class="buttons" v-if="othername != username">
				<button class="unban-button" v-if="banning" :onclick="() => unbanThis().then()">Unban</button>
				<button class="ban-button" v-if="!banning" :onclick="() => banThis().then()">Ban</button>
				<button class="unfollow-button" v-if="following" :onclick="() => unfollowThis().then()">Unfollow</button>
				<button class="follow-button" v-if="!following" :onclick="() => followThis().then()">Follow</button>
			</div>
			<div><span v-if="othername != username"><b>{{ othername }}</b>'s</span><span v-if="othername == username"><b>Your</b></span> followers:</div>
			<div class="follow-container">
				<router-link class="follow-element" v-for="element in followersResults"  v-bind:to="'/users/' +element.name">{{ element.name }}</router-link>
				<router-link v-if="following" class="follow-element" :to="'/users/' + username">{{ username }}</router-link>
				<button class="follow-element-end" :onclick="() => refreshData(false, true, false)">More followers!</button>
			</div>
			<div><span v-if="othername != username"><b>{{ othername }}</b>'s</span><span v-if="othername == username"><b>Your</b></span> following:</div>
			<div class="follow-container">
				<router-link class="follow-element" v-for="element in followingResults"  v-bind:to="'/users/' +element.name">{{ element.name }}</router-link>
				<button class="follow-element-end" :onclick="() => refreshData(false, false, true)">More followings!</button>
			</div>
			<div v-if="othername==username">
				Modify your username:<br/>
				<input class="new-username" type="text" v-on:keyup.enter="newUsername" v-model="input" placeholder="New username" maxlength="32"/>
			</div>
		</div>
		<div class="posts-holder">
            <PostCard v-for="elem in photosResults" v-bind:imgId="elem" v-bind:del="()=>delPost(elem)"></PostCard>
            <button class="posts-more" :onclick="() => refreshData(true, false, false)">Show more posts!</button>
        </div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>