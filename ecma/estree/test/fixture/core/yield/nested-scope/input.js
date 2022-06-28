async function* generate(stream, boundary, options) {
  outer: for await (const chunk of stream) {
    is_eager ? yield tmp : payloads.push(tmp);
  }
}
