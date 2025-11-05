// https://github.com/adyliu/china_area
// id, parent, name, pinyin, post, address

const py = require(`pinyin`);
const raw = require(`fs`).readFileSync(`./cn-raw.csv`, `utf8`)
x = raw.split(`\n`).map(l => l.split`,`).map(a => ({
  id: a[0],
  name: a[1],
  parent: parseInt(a[3]),
}))

map = Object.fromEntries(x.map(a => [a.id, a]))


res = x.map(a => {
  let address = a.name
  let parent = map[a.parent]
  while (parent) {
    address = parent.name + address
    parent = map[parent.parent]
  }
  const pinyin = py.pinyin(a.name, { style: "normal" }).flat().join(` `)
  return {
    id: a.id,
    parent: a.parent,
    name: a.name,
    pinyin,
    address
  }
})

const csv = res.map(a => `${a.id},${a.parent || ''},${a.name},${a.pinyin},${a.address || ''},`).join('\n')

require(`fs`).writeFileSync(`../data/cn.csv`, csv)