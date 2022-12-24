<script>
import PostPreview from '../components/PostPreview.vue';

export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
        };
    },
    components: { PostPreview }
}
</script>

<template>
	<div class="posts-holder"></div>
    <ol class="no-bullets">
        <il>
            <PostPreview></PostPreview>
        </il>
    </ol>
</template>