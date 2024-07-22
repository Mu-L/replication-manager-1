import { Box, Flex, Image, Button, Spacer, Text, HStack, useColorMode, IconButton } from '@chakra-ui/react'
import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { logout } from '../redux/authSlice'
import ThemeIcon from './ThemeIcon'
import RefreshCounter from './RefreshCounter'
import { isAuthorized } from '../utility/common'
import { Link } from 'react-router-dom'
import { useTheme } from '@emotion/react'
import { clearCluster } from '../redux/clusterSlice'
import AlertBadge from './AlertBadge'
import AlertModal from './Modals/AlertModal'
import { FaPowerOff } from 'react-icons/fa'
import ConfirmModal from './Modals/ConfirmModal'

function Navbar({ username }) {
  const dispatch = useDispatch()
  const { colorMode } = useColorMode()
  const [alertModalType, setAlertModalType] = useState('')
  const [isLogoutModalOpen, setIsLogoutModalOpen] = useState(false)
  const {
    common: { isMobile, isTablet, isDesktop },
    cluster: { clusterAlerts }
  } = useSelector((state) => state)

  const currentTheme = useTheme()

  const styles = {
    navbarContainer: {
      boxShadow: colorMode === 'dark' ? 'none' : '0px -1px 8px #BFC1CB',
      position: 'fixed',
      zIndex: 2,
      width: '100%',
      padding: '4px',
      background: colorMode === 'light' ? currentTheme.colors.primary.light : currentTheme.colors.primary.dark
    },
    logo: {
      bg: '#eff2fe',
      borderRadius: '4px'
    },
    alerts: {
      gap: '4'
    }
  }

  const openAlertModal = (type) => {
    setAlertModalType(type)
  }
  const closeAlertModal = (type) => {
    setAlertModalType('')
  }

  const openLogoutModal = () => {
    setIsLogoutModalOpen(true)
  }
  const closeLogoutModal = () => {
    setIsLogoutModalOpen(false)
  }

  const handleLogout = () => {
    dispatch(logout())
    dispatch(clearCluster())
  }
  return (
    <>
      <Flex as='nav' sx={styles.navbarContainer} gap='2' align='center'>
        <Link to='/'>
          <Image
            loading='lazy'
            height='50px'
            width={isMobile ? '180px' : 'fit-content'}
            sx={styles.logo}
            objectFit='contain'
            src='/images/logo.png'
            alt='Replication
           Manager'
          />
        </Link>
        <Spacer />

        {isAuthorized() && isDesktop && <RefreshCounter />}

        <Spacer />
        <HStack spacing='4'>
          {isAuthorized() && (
            <Flex sx={styles.alerts}>
              <AlertBadge
                isBlocking={true}
                text='Blockers'
                count={clusterAlerts?.errors?.length || 0}
                onClick={() => openAlertModal('error')}
                showText={!isMobile}
              />
              <AlertBadge
                text='Warnings'
                count={clusterAlerts?.warnings?.length || 0}
                onClick={() => openAlertModal('warning')}
                showText={!isMobile}
              />
            </Flex>
          )}

          {isAuthorized() && (
            <>
              {username && isDesktop && <Text>{`Welcome, ${username}`}</Text>}
              {isMobile ? (
                <IconButton
                  onClick={openLogoutModal}
                  variant='filled'
                  border='none'
                  size='md'
                  icon={<FaPowerOff fontSize='1.5rem' fill='blue.200' />}
                />
              ) : (
                <Button type='button' size={{ base: 'sm' }} onClick={openLogoutModal}>
                  Logout
                </Button>
              )}
            </>
          )}

          <ThemeIcon />
        </HStack>
      </Flex>
      {isAuthorized() && !isDesktop && (
        <Box mx='auto' p='16px' marginTop='80px'>
          <RefreshCounter />
        </Box>
      )}
      {alertModalType && (
        <AlertModal type={alertModalType} isOpen={alertModalType.length !== 0} closeModal={closeAlertModal} />
      )}
      {isLogoutModalOpen && (
        <ConfirmModal
          onConfirmClick={handleLogout}
          closeModal={closeLogoutModal}
          isOpen={isLogoutModalOpen}
          title={'Are you sure you want to log out?'}
        />
      )}
    </>
  )
}

export default Navbar
