<template>
  <div>
    <p>{{count}}</p>
    <p>{{year}}</p>
    <ul>
<!--      渲染列表元素 将li的DOM key属性设置为item的id 展示每一个name-->
      <li v-for="item in list" :key="item.id">{{item.name}}</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from "vue";

const count = ref(39)

// 使用TS中泛型处理数据
const year = ref<string | number>('9393')

type ListItem = {
  id:number
  name:string
}

type Songs = {
  name: string
  producer: string
  // 使用可选运算符对字段进行处理
  language?: string
}
// 创建自定义的数组(空)
const list = ref<ListItem[]>([])

// 直接向数组中插入元素
list.value = [{
  id: 39,
  name: 'miku'
}]

// reactive对象会自动推导字段的类型
// const songs = reactive({
//   name: '' ,
//   producer: ''
// })
// 不能再直接额外插入新的字段了
// songs.language = 'jp'

// 使用泛型创建对象
const songs = reactive<Songs>({
  name: '39',
  producer: 'sskr&deco27',
  language: 'JP'
})
</script>

<style scoped>

</style>