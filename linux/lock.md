中断屏蔽:
(可以保证正在执行的内核执行路径不被中断处理程序抢占,由于Linux内核的进程调度都依赖中断来实现，内核抢占进程之间的竞态就不存在了)

使用方法: 
	local_irq_disable()  //屏蔽中断                                   说明：local_irq_disable()和local_irq_enable()都只能禁止和使能本CPU内的中断
	….                                                                               并不能解决SMP多CPU引发的竞争。
	critical section  //临界区
	….

    local_irq_enable()  //开中断 

	与local_irq_disable()不同，local_irq_save(flags)除了进行禁止中断操作以外，还保证目前CPU的中断位信息，local_irq_restore(flags)进行相反的操作。

自旋锁:
正如其名,CPU上将要执行的代码将会执行一个测试并设置某个内存变量的原子操作,若测试结果表明锁已经空闲,则程序获得这个自旋
	锁继续运行;若仍被占用,则程序将在一个小的循环内重复测试这个"测试并设置"的操作.这就是自旋。
	使用方法:1)spinlock_t spin;  //定义自旋锁
	2)spin_lock_init(lock);  //初始化自旋锁
	3)spin_lock(lock);  //成功获得自旋锁立即返回，否则自旋在那里直到该自旋锁的保持者释放
	spin_trylock(lock); //成功获得自旋锁立即返回真，否则返回假，而不是像上一个那样"在原地打转"
	4)spin_unlock(lock);//释放自旋锁
	自旋锁一般像下边这样使用：
	spinlock_t lock;
	spin_lock_init(&lock);
	spin_lock (&lock);
	....//临界区
	spin_unlock(&lock);
	还记的前边说的第一招：中断屏蔽中致命的弱点步，自旋锁就是针对SMP或单个CPU但内核可抢占的情况，对于但CPU和内核不可抢占的系统，自旋锁退化为空操作。还有就是自旋锁解决了临界区不受别的CPU和本CPU内的抢占进程打扰，但是得到锁的代码路径在执行临界区的时候还可能受到中断和底半部的影响。

利用spin_lock()/spin_unlock()作为自旋锁的基础，将它们和关中断local_irq_disable()/开中断local_irq_enable(),关底半部local_bh_disable()/开底半部local_bh_enable(),关中断并保存状态字local_irq_save()/开中断并恢复状态local_irq_restore()结合就完成了整套自旋锁机制。
	spin_lock_irq() = spin_lock() + local_irq_disable()
	spin_unlock_irq = spin_unlock() + local_irq_enable()
	spin_lock_irqsave() = spin_unlock() + local_irq_save()
	spin_unlock_irqrestore() = spin_unlock() + local_irq_restore()
	spin_lock_bh() = spin_lock() + local_bh_disable()
spin_unlock_bh() = spin_unlock() +local_bh_enable()


读写自旋锁:
	它保留了自锁的概念，但是它规定在读方面同时可以有多个读单元，在写方面，只能最多有一个写进程。当然，读和写也不能同时进行。
	使用方法：
	1)初始化读写锁的方法。
		rwlock_t x;//静态初始化 rwlock_t x=RW_LOCK_UNLOCKED;//动态初始化 rwlock_init(&x);
	2)最基本的读写函数。
		void read_lock(rwlock_t *lock);//使用该宏获得读写锁，如果不能获得，它将自旋，直到获得该读写锁
		void read_unlock(rwlock_t *lock);//使用该宏来释放读写锁lock
		void write_lock(rwlock_t *lock);//使用该宏获得获得读写锁，如果不能获得，它将自旋，直到获得该读写锁
		void write_unlock(rwlock_t *lock);//使用该宏来释放读写锁lock
	3)和自旋锁中的spin_trylock(lock),读写锁中分别为读写提供了尝试获取锁，并立即返回的函数，如果获得，就立即返回真，否则返回假：
		read_trylock(lock)和write_lock(lock);

	4)硬中断安全的读写锁函数：
		read_lock_irq(lock);//读者获取读写锁，并禁止本地中断
		read_unlock_irq(lock);//读者释放读写锁，并使能本地中断
		write_lock_irq(lock);//写者获取读写锁，并禁止本地中断
		write_unlock_irq(lock);//写者释放读写锁，并使能本地中断
		read_lock_irqsave(lock, flags);//读者获取读写锁，同时保存中断标志，并禁止本地中断
		read_unlock_irqrestores(lock,flags);//读者释放读写锁，同时恢复中断标志，并使能本地中断
		write_lock_irqsave(lock,flags);//写者获取读写锁，同时保存中断标志，并禁止本地中断
		write_unlock_irqstore(lock,flags);写者释放读写锁，同时恢复中断标志，并使能本地中断
	5)软中断安全的读写函数：
		read_lock_bh(lock);//读者获取读写锁，并禁止本地软中断
		read_unlock_bh(lock);//读者释放读写锁，并使能本地软中断
		write_lock_bh(lock);//写者获取读写锁，并禁止本地软中断
		write_unlock_bh(lock);//写者释放读写锁，并使能本地软中断

信号量(信号量其实和自旋锁是一样的，就是有一点不同:当获取不到信号量时，进程不会原地打转而是进入休眠等待状态)
	Linux系统中与信号量相关的操作主要有一下4种：
	1)定义信号量    struct semaphore sem;
	2)初始化信号量   
		void sema_init (struct semphore *sem, int val);    //设置sem为val
		void init_MUTEX(struct semaphore *sem);    //初始化一个用户互斥的信号量sem设置为1
		void init_MUTEX_LOCKED(struct semaphore *sem);    //初始化一个用户互斥的信号量sem设置为0
		DECLARE_MUTEX(name);     //该宏定义信号量name并初始化1
		DECLARE_MUTEX_LOCKED(name);    //该宏定义信号量name并初始化0
	3)获得信号量
		void down(struct semaphore *sem);    //该函数用于获取信号量sem，会导致睡眠，不能被信号打断,所以不能在中断上下文使用。
		int down_interruptible(struct semaphore *sem);    //因其进入睡眠状态的进程能被信号打断，信号也会导致该函数返回，这是返回非0。
		int down_trylock(struct semaphore *sem);//尝试获得信号量sem，如果能够获得，就获得并返回0，否则返回非0,不会导致调用者睡眠，可以在中断上下文使用
		一般这样使用
			if(down_interruptible(&sem))
			{
				return  - ERESTARTSYS;
			}

	4)释放信号量
		void up(struct semaphore *sem);    //释放信号量sem，唤醒等待者
	信号量一般这样被使用，如下所示：
		//定义信号量
		DECLARE_MUTEX(mount_sem);
		down(&mount_sem);//获取信号量，保护临界区
			…
		critical section //临界区
			…
		up(&mount_sem);
