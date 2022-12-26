<script>
import PostCard from './PostCard.vue';

export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
        };
    },
    components: { PostCard }
}
</script>

<template>
    <div class="posts-holder">
        <PostCard imgId=1></PostCard>
        <PostCard></PostCard>
    </div>
</template>