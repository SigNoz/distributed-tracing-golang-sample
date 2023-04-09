// import require from './require'
// require("dotenv").config()

// Components
const Title = (props) => <p>{props.label}</p>
const Response = (props) => <pre>{JSON.stringify(props.result, null, 4)}</pre>
const ProductDD = (props) => {
  return (
    <select onChange={props.handleProductChange}>
      <option selected>Choose product</option>
      <option value={999}>MacBook Air</option>
      <option value={1100}>MacBook Pro</option>
      <option value={599}>Mac Mini</option>
    </select>
  )
}

const App = () => {
  let [user, setUser] = React.useState()
  let [payment, setPayment] = React.useState()
  let [order, setOrder] = React.useState()
  let [product, setProduct] = React.useState()
  let [signozUrl, setSignozUrl] = React.useState(DEFAULT_SIGNOZ_URL)
  let createUser = async () => {
    try {
      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_name: 'Naman',
          account: 'savings',
        }),
      }
      let userServiceUrl='/users';
      if (typeof USER_PORT !== 'undefined') {
        userServiceUrl = `http://localhost:${USER_PORT}` + userServiceUrl
      }
      let response = await fetch(userServiceUrl, requestOptions)
      let result = await response.json()
      setUser(result)
      console.log(result)
    } catch (error) {
      setUser(error.message)
      console.log('err', error.message)
    }
  }

  let handleProduct = (e) => {
    let { options, value } = e.target
    setProduct({
      product_name: options[options.selectedIndex].text,
      price: value,
    })
  }

  let transferFund = async () => {
    console.log('transferFund', user)
    const requestOptions = {
      method: 'PUT',
	  headers: { 'Content-Type': 'application/json' },
	  body: JSON.stringify({
	    amount: 10000,
	  }),
    }

    let paymentServiceUrl='/payments/transfer/id/${user.id}';
    if (typeof PAYMENT_PORT !== 'undefined') {
      paymentServiceUrl = `http://localhost:${PAYMENT_PORT}` + paymentServiceUrl
    }
    let response = await fetch(
      paymentServiceUrl,
      requestOptions
    )
    let result = await response.json()
    setPayment(result)
    console.log(result)
  }

  let placeOrder = async () => {
    console.log('placeOrder', product)
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: user.id,
        product_name: product.product_name,
        price: parseInt(product.price),
      }),
    }

    let orderServiceUrl='/orders';
    if (typeof ORDER_PORT !== 'undefined') {
      orderServiceUrl = `http://localhost:${ORDER_PORT}` + orderServiceUrl
    }
    let response = await fetch(orderServiceUrl, requestOptions)
    let result = await response.json()
    setOrder(result)
    console.log(result)
  }

  let handleReset = () => {
    setUser()
    setPayment()
    setOrder()
    setProduct()
  }

  let setUrl = (e) => setSignozUrl(e.target.value)
  let openSignoz = () => window.open(signozUrl)
  return (
    <div>
      <button onClick={handleReset}>Reset Actions</button>
      <button onClick={openSignoz} style={{ marginLeft: 10, marginRight: 8 }}>
        Open Signoz
      </button>
      <input
        placeholder="Type default signoz url"
        type="text"
        onChange={setUrl}
        value={signozUrl}
        style={{ width: 200 }}
      ></input>
      <p style={{ color: 'red' }}>Make sure signoz is running on {signozUrl}</p>
      <Title label="1. User Creation"></Title>
      <button onClick={createUser}>Create User</button>
      <Response result={user} />
      {user && (
        <div>
          <Title label="2. Transfer amount"></Title>
          <button onClick={transferFund}>Transfer Fund</button>
          <Response result={payment} />
        </div>
      )}
      {payment && (
        <div>
          <Title label="3. Place order"></Title>
          <ProductDD handleProductChange={handleProduct} />
          {product && <button onClick={placeOrder}>Place Order</button>}
          {order && (
            <div>
              <Response result={order} />
              <h3>Order Placed!</h3>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

ReactDOM.render(<App />, document.getElementById('app'))

