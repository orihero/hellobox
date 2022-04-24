import React, { useEffect, useState } from "react";
import { requests } from "../../../../api/requests";
import {
  Activated, Container, Icon,
  IconBox, ProdoctsBox, ProdoctsContainer, Sold, Table,
  Td, Tdbox, Tr, Wrapper
} from "./style";

function OrderPage() {
  const [orderIndex, setOrderIndex] = useState(-1);
  const onClickDrop = (i) => {
    setOrderIndex(i);
  };

  const [orders, setOrders] = useState([]);

  const effect = async () => {
    try {
      let res = await requests.order.getOrders();
      setOrders(res.data);
    } catch (error) {
      alert("Error in getting orders");
    }
  };

  useEffect(() => {
    effect();
  }, []);

  return (
    <Container>
      {orders.map((e, i) => {
        return (
          <>
            <ProdoctsBox>
              <Table>
                <Tr>
                  <Td>{e.id}</Td>
                  <Td>{e.user.first_name}</Td>
                  <Td>{e.user.phone}</Td>
                  <Td>
                    <Tdbox>
                      <Sold>
                        Sold: {new Date(e.created_at).toLocaleDateString()}
                      </Sold>
                    </Tdbox>
                  </Td>
                </Tr>
              </Table>
              <Icon>
                <IconBox
                  onClick={() => {
                    console.log({ i });
                    onClickDrop(i);
                  }}
                >
                  {orderIndex === i ? (
                    <i class="bx bx-chevron-up"></i>
                  ) : (
                    <i class="bx bx-chevron-down"></i>
                  )}
                </IconBox>
              </Icon>
            </ProdoctsBox>
            <Wrapper
              style={{
                height:
                  i === orderIndex
                    ? `${e.cart.products.length * 100}px`
                    : "0px",
              }}
            >
              {e.cart.products.map((el) => {
                return (
                  <ProdoctsContainer>
                    {/* Product component */}
                    <Table>
                      <Tr>
                        <Td>{el.product.name}</Td>
                        <Td>{el.product.count}</Td>
                        <Td>{el.product.price}</Td>
                        <Td>{el.token}</Td>
                        <Td>{el.utilized ? "Activated ✅" : "Active ❌"}</Td>
                        <Td>
                          <Tdbox>
                            <Sold>Sold: 10/20/2021</Sold>
                            <Activated>Activated: 20/10/2021</Activated>
                          </Tdbox>
                        </Td>
                      </Tr>
                    </Table>
                  </ProdoctsContainer>
                );
              })}
            </Wrapper>
          </>
        );
      })}
    </Container>
  );
}

export default OrderPage;
